package plugin

import (
	"fmt"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"sync"
)

//GetOAuthConfig returns the oauth config
func (p *Plugin) GetOAuthConfig() *oauth2.Config {

	config := p.getConfiguration()
	apiConfig := p.API.GetConfig()

	redirectURL := fmt.Sprintf("%s/plugins/mattermost-plugin-servicenow/oauth/complete", *apiConfig.ServiceSettings.SiteURL)
	return &oauth2.Config{
		ClientID:     config.ApplicationID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth_auth.do", config.ServiceNowURL),
			TokenURL: fmt.Sprintf("%s/oauth_token.do", config.ServiceNowURL),
		},
		RedirectURL: redirectURL,
	}
}

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

	if err := p.API.RegisterCommand(getCommand()); err != nil {
		return errors.Wrap(err, "failed to register servicenow command")
	}

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
