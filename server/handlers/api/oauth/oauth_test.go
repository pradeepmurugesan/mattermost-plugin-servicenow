package oauth

import (
	"bytes"
	"fmt"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	mock2 "github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOauth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oauth")
}

func setUpTestOauthServer() *httptest.Server {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/oauth_token.do" {
			Fail(fmt.Sprintf("Unexpected exchange request URL %q", r.URL))
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Fail(fmt.Sprintf("Failed reading request body: %s.", err))
		}
		if string(body) != "code=hello&grant_type=authorization_code&redirect_uri=REDIRECT_URL" {
			Fail(fmt.Sprintf("Unexpected exchange payload; got %q", body))
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=some-token&scope=user&token_type=bearer"))
	}))

	return ts
}

func setUpServiceNowMockServer() *httptest.Server {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/api/me" {
			Fail(fmt.Sprintf("Unexpected request URL %q", r.URL))
		}

		auth := r.Header.Get("Authorization")

		if auth != "Bearer some-token" {
			Fail("Unauthorized")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"result\":{\"id\":\"681633be04be441\",\"name\":\"admin\",\"display_name\":\"System Administrator\",\"first_name\":\"System\",\"last_name\":\"Administrator\"}}"))
	}))

	return ts
}
func getOauthConfig(url string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth_auth.do", url),
			TokenURL: fmt.Sprintf("%s/oauth_token.do", url),
		},
		RedirectURL: "REDIRECT_URL",
	}
}

var _ = Describe("Oauth", func() {

	var (
		pluginAPIMock   *plugintest.API
		pluginContext   *models.PluginContext
		mockHTTPRequest *http.Request
		oauthConfig     *oauth2.Config
	)

	BeforeEach(func() {
		pluginAPIMock = &plugintest.API{}
		oauthConfig = getOauthConfig("http://localhost")
		pluginContext = &models.PluginContext{BotUserID: "some-userId", OauthConfig: *oauthConfig, API: pluginAPIMock}

	})

	Describe("Authorize", func() {

		BeforeEach(func() {
			mockHTTPRequest = httptest.NewRequest("GET", "/oauth/connect", bytes.NewReader([]byte("{}")))
		})

		It("should throw an error in case the header doesn't contain Mattermost-User-ID", func() {

			response := httptest.NewRecorder()

			Authorize(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusUnauthorized))
			Expect(response.Body.String()).To(Equal("Not authorized\n"))

		})

		It("should redirect to the authorize url provided in the oauthConfig", func() {

			mockHTTPRequest.Header.Set("Mattermost-User-ID", "some-id")
			response := httptest.NewRecorder()
			pluginAPIMock.On("KVSet", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
				Return(nil)

			Authorize(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusPermanentRedirect))

		})
	})

	Describe("CompleteAuthentication", func() {

		BeforeEach(func() {
			mockHTTPRequest = httptest.NewRequest("GET", "/oauth/complete", bytes.NewReader([]byte("{}")))
		})

		It("should throw an error in case the header doesn't contain Mattermost-User-ID", func() {

			response := httptest.NewRecorder()

			CompleteAuthentication(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusUnauthorized))
			Expect(response.Body.String()).To(Equal("Not authorized\n"))

		})

		It("should throw an error in case the query param code is missing", func() {

			mockHTTPRequest.Header.Set("Mattermost-User-ID", "some-id")
			response := httptest.NewRecorder()

			CompleteAuthentication(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.Body.String()).To(Equal("missing authorization code\n"))

		})

		It("should throw an error in case the stored state is not found", func() {

			mockHTTPRequest = httptest.NewRequest("GET", "/oauth/complete?code=hello&state=123", bytes.NewReader([]byte("{}")))
			mockHTTPRequest.Header.Set("Mattermost-User-ID", "some-id")
			response := httptest.NewRecorder()
			pluginAPIMock.On("KVGet", "123").Return(nil, &model.AppError{Message: "Error from mock"})

			CompleteAuthentication(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.Body.String()).To(Equal("missing stored state\n"))

		})

		It("should throw an error in case the stored state is different", func() {

			mockHTTPRequest = httptest.NewRequest("GET", "/oauth/complete?code=hello&state=123", bytes.NewReader([]byte("{}")))
			mockHTTPRequest.Header.Set("Mattermost-User-ID", "some-id")
			response := httptest.NewRecorder()
			pluginAPIMock.On("KVGet", "123").Return([]byte("1234"), nil)

			CompleteAuthentication(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.Body.String()).To(Equal("invalid state\n"))

		})

		It("should throw an error in case the stored userId and header User Id are different", func() {

			mockHTTPRequest = httptest.NewRequest("GET", "/oauth/complete?code=hello&state=123_123", bytes.NewReader([]byte("{}")))
			mockHTTPRequest.Header.Set("Mattermost-User-ID", "some-id")
			response := httptest.NewRecorder()
			pluginAPIMock.On("KVGet", "123_123").Return([]byte("123_123"), nil)
			pluginAPIMock.On("KVDelete", "123_123").Return(nil)

			CompleteAuthentication(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusUnauthorized))
			Expect(response.Body.String()).To(Equal("Not authorized, incorrect user\n"))

		})

		It("should exchange the code for token", func() {

			ts := setUpTestOauthServer()
			sn := setUpServiceNowMockServer()
			oauthConfig = getOauthConfig(ts.URL)
			pluginContext = &models.PluginContext{BotUserID: "some-userId", OauthConfig: *oauthConfig, API: pluginAPIMock, UserInfoEndpoint: fmt.Sprintf("%s/api/me", sn.URL)}
			pluginAPIMock.On("KVSet", models.NowTokenKeyPrefix+"some-id", mock2.AnythingOfType("[]uint8")).Return(nil)
			pluginAPIMock.On("KVSet", models.NowUserIDPrefix+"681633be04be441", mock2.AnythingOfType("[]uint8")).Return(nil)

			mockHTTPRequest = httptest.NewRequest("GET", "/oauth/complete?code=hello&state=123_some-id", bytes.NewReader([]byte("{}")))
			mockHTTPRequest.Header.Set("Mattermost-User-ID", "some-id")
			response := httptest.NewRecorder()
			pluginAPIMock.On("KVGet", "123_some-id").Return([]byte("123_some-id"), nil)
			pluginAPIMock.On("KVDelete", "123_some-id").Return(nil)

			CompleteAuthentication(response, mockHTTPRequest, pluginContext)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.Body.String()).To(ContainSubstring("Completed connecting to ServiceNow."))

		})
	})
})
