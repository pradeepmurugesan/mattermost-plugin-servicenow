package models

import (
	"github.com/mattermost/mattermost-server/plugin"
	"golang.org/x/oauth2"
)

// PluginContext context of the plugin
type PluginContext struct {
	API         plugin.API
	BotUserID   string
	OauthConfig oauth2.Config
}

// SetAPI sets the api
func (p *PluginContext) SetAPI(api plugin.API) {
	p.API = api
}
