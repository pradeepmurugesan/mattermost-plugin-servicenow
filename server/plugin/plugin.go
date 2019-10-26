package plugin

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	BotUserID string
	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex
	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

//OnActivate hook to be called after the plugin activation
func (p *Plugin) OnActivate() error {

	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "service-now",
		DisplayName: "Service Now",
		Description: "Created by the Service Now plugin.",
	})

	if err != nil {
		return errors.Wrap(err, "failed to ensure service-now bot")
	}

	p.BotUserID = botID

	return nil

}
