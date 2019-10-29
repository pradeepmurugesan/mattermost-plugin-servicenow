package stream

import (
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestStreamCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stream Suite")
}

var _ = Describe("Stream", func() {

	var (
		pluginAPIMock *plugintest.API
		pluginContext *models.PluginContext
	)

	BeforeEach(func() {
		pluginAPIMock = &plugintest.API{}
		pluginContext = &models.PluginContext{BotUserID: "some-userId"}
		pluginContext.SetAPI(pluginAPIMock)
	})
	const MockChannelID = "mock-channel-id"

	It("should throw the error in case of KVGet failure", func() {

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, &model.AppError{Message: "Error from the mock"})

		_, err := Execute(pluginContext, &model.CommandArgs{})

		Expect(err.Message).To(Equal("Error from the mock"))
	})

	It("should throw error in case stream is already subscribed", func() {

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return([]byte("abcd123"), nil)

		_, err := Execute(pluginContext, &model.CommandArgs{})

		Expect(err.Message).To(Equal("Already subscribed to service now stream"))
	})

	It("should throw error in case GetChannel fails", func() {

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, nil)
		pluginAPIMock.On("GetChannel", MockChannelID).
			Return(nil, &model.AppError{Message: "Error from the GetChannel mock"})

		_, err := Execute(pluginContext, &model.CommandArgs{ChannelId: MockChannelID})

		Expect(err.Message).To(Equal("Error from the GetChannel mock"))
	})

	It("should subscribe to the channel from the command args", func() {

		pluginAPIMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, nil)
		pluginAPIMock.On("GetChannel", MockChannelID).
			Return(&model.Channel{Id: MockChannelID, Name: "mock-channel-name"}, nil)
		pluginAPIMock.On("KVSet", models.ServiceNowSubscribedChannel, []byte(MockChannelID)).
			Return(nil, nil)

		result, err := Execute(pluginContext, &model.CommandArgs{ChannelId: MockChannelID})

		Expect(err).To(BeNil())
		Expect(result).To(Equal("subscribed for the now incident updates to the channel: mock-channel-name"))
	})
})
