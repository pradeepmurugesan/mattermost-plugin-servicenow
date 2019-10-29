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
		pluginApiMock *plugintest.API
		pluginContext *models.PluginContext
	)

	BeforeEach(func() {
		pluginApiMock = &plugintest.API{}
		pluginContext = &models.PluginContext{BotUserId: "some-userId"}
		pluginContext.SetApi(pluginApiMock)
	})
	const MockChannelId = "mock-channel-id"

	It("should throw the error in case of KVGet failure", func() {

		pluginApiMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, &model.AppError{Message: "Error from the mock"})

		_, err := Execute(pluginContext, &model.CommandArgs{})

		Expect(err.Message).To(Equal("Error from the mock"))
	})

	It("should throw error in case stream is already subscribed", func() {

		pluginApiMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return([]byte("abcd123"), nil)

		_, err := Execute(pluginContext, &model.CommandArgs{})

		Expect(err.Message).To(Equal("Already subscribed to service now stream"))
	})

	It("should throw error in case GetChannel fails", func() {

		pluginApiMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, nil)
		pluginApiMock.On("GetChannel", MockChannelId).
			Return(nil, &model.AppError{Message: "Error from the GetChannel mock"})

		_, err := Execute(pluginContext, &model.CommandArgs{ChannelId: MockChannelId})

		Expect(err.Message).To(Equal("Error from the GetChannel mock"))
	})

	It("should subscribe to the channel from the command args", func() {

		pluginApiMock.On("KVGet", models.ServiceNowSubscribedChannel).
			Return(nil, nil)
		pluginApiMock.On("GetChannel", MockChannelId).
			Return(&model.Channel{Id: MockChannelId, Name: "mock-channel-name"}, nil)
		pluginApiMock.On("KVSet", models.ServiceNowSubscribedChannel, []byte(MockChannelId)).
			Return(nil, nil)

		result, err := Execute(pluginContext, &model.CommandArgs{ChannelId: MockChannelId})

		Expect(err).To(BeNil())
		Expect(result).To(Equal("subscribed for the now incident updates to the channel: mock-channel-name"))
	})
})
