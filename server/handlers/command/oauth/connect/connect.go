package connect

import (
	"fmt"
	"github.com/mattermost/mattermost-server/model"
)

//Execute executes the command
func Execute(config *model.Config) string {
	link := fmt.Sprintf("%s/plugins/mattermost-plugin-servicenow/oauth/connect", *config.ServiceSettings.SiteURL)
	return fmt.Sprintf("[click here](%s) to connect with ServiceNow", link)
}
