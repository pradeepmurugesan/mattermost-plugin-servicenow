package oauth

import (
	"bytes"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOauth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oauth")
}


var _ = Describe("Oauth", func() {

	var (
		pluginAPIMock  *plugintest.API
		pluginContext   *models.PluginContext
		mockHTTPRequest *http.Request
		oauthConfig *oauth2.Config
	)

	BeforeEach(func() {
		pluginAPIMock = &plugintest.API{}
		oauthConfig = &oauth2.Config{
			ClientID:     "client",
			ClientSecret: "secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://localhost/oauth_auth.do",
				TokenURL: "http://localhost/oauth_token.do",
			},
			RedirectURL: "http://localhost/oauth_redirect.do",
		}
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
			pluginAPIMock.On("KVSet", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8") ).
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
	})
})
