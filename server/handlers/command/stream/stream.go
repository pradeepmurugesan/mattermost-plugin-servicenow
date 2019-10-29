package stream

import (
	"fmt"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/model"
)

func Execute(ctx *models.PluginContext, args *model.CommandArgs) (string, *models.Error) {

	subscribedChannelId, err := ctx.API.KVGet(models.ServiceNowSubscribedChannel)

	if err != nil {
		return "", &models.Error{Message: err.Message}
	}


	if subscribedChannelId != nil {
		return "", &models.Error{Message: "Already subscribed to service now stream"}
	}

	channel, err := ctx.API.GetChannel(args.ChannelId)

	if err != nil {
		return "", &models.Error{Message: err.Message}
	}

	if err = ctx.API.KVSet(models.ServiceNowSubscribedChannel, []byte(args.ChannelId)); err != nil {
		return "", &models.Error{Message: err.Message}
	}

	return fmt.Sprintf("subscribed for the now incident updates to the channel: %s", channel.Name), nil
}
