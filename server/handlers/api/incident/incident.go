package incident

import (
	"encoding/json"
	"fmt"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-plugin-servicenow/server/template"
	"github.com/mattermost/mattermost-server/model"
	"net/http"
)

// PublishIncident publish to channel if subscribed
func PublishIncident(writer http.ResponseWriter, request *http.Request, ctx models.PluginContext) *models.Error {

	subscribedChannelID, err := ctx.API.KVGet(models.ServiceNowSubscribedChannel)

	if err != nil {
		return models.NewError(err.Message)
	}

	if subscribedChannelID == nil {
		return models.NewError("ServiceNow stream subscribedChannelID is nil")
	}

	d := json.NewDecoder(request.Body)
	d.DisallowUnknownFields()
	var incidentObj models.Incident

	if err := d.Decode(&incidentObj); err != nil {
		return models.NewError(fmt.Sprintf("Couldn't decode: %v", err))
	}

	postMessage, e := template.RenderTemplate("incident", incidentObj)
	if e != nil {
		return models.NewError(fmt.Sprintf("Couldn't render incident template: %s", e.Error()))
	}

	post := &model.Post{
		UserId:    ctx.BotUserID,
		ChannelId: string(subscribedChannelID),
		Message:   postMessage,
	}
	_, err = ctx.API.CreatePost(post)

	if err != nil {
		return models.NewError(fmt.Sprintf("Couldn't create incident post: %s", err.Message))
	}

	return nil
}
