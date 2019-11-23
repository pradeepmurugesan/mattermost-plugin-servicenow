package incident

import (
	"bytes"
	"encoding/json"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-plugin-servicenow/server/template"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublishIncident(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Publish Incident")
}

var _ = Describe("Incident", func() {
	var (
		pluginAPIMock  *plugintest.API
		pluginContext  *models.PluginContext
		mockHTTPRequst *http.Request
	)

	BeforeEach(func() {
		pluginAPIMock = &plugintest.API{}
		pluginContext = &models.PluginContext{BotUserID: "some-userId"}
		pluginContext.SetAPI(pluginAPIMock)
		mockHTTPRequst = httptest.NewRequest("POST", "/incident", bytes.NewReader([]byte("{}")))
	})

	It("should throw the error in case of KVGet failure", func() {

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, &model.AppError{Message: "Error from the mock"})

		err := PublishIncident(httptest.NewRecorder(), mockHTTPRequst, *pluginContext)

		Expect(err.Message).To(Equal("Error from the mock"))
	})

	It("should throw an error in case of subscribedChannelID is nil", func() {

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, nil)

		err := PublishIncident(httptest.NewRecorder(), mockHTTPRequst, *pluginContext)

		Expect(err.Message).To(Equal("ServiceNow stream subscribedChannelID is nil"))
	})

	It("should throw an error in case of invalid json", func() {

		mockHTTPRequst = httptest.NewRequest("POST", "/incident", bytes.NewReader([]byte("invalid json")))

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return([]byte("some-id"), nil)

		err := PublishIncident(httptest.NewRecorder(), mockHTTPRequst, *pluginContext)

		Expect(err.Message).To(Equal("Couldn't decode: invalid character 'i' looking for beginning of value"))
	})

	It("should call the createPost function with the right parameters", func() {

		var incident = models.Incident{
			SysCreatedBy:     "some-user",
			ShortDescription: "description of the ticket",
			CreatedByID:      "some-id",
			Priority:         "1",
			Impact:           "5",
		}

		incidentJSON, _ := json.Marshal(incident)
		expectedMessage, _ := template.RenderTemplate("incident", incident)

		mockHTTPRequst = httptest.NewRequest("POST", "/incident", bytes.NewReader(incidentJSON))
		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return([]byte("some-id"), nil)
		pluginAPIMock.On("CreatePost", mock.MatchedBy(func(post *model.Post) bool {
			return post.UserId == "some-userId" && post.Message == expectedMessage
		})).Return(&model.Post{}, nil)

		err := PublishIncident(httptest.NewRecorder(), mockHTTPRequst, *pluginContext)

		Expect(err).To(BeNil())

	})
})
