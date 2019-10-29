package models

import (
	"github.com/mattermost/mattermost-server/plugin"
)

// PluginContext context of the plugin
type PluginContext struct {
	API       plugin.API
	BotUserID string
}

// SetAPI sets the api
func (p *PluginContext) SetAPI(api plugin.API) {
	p.API = api
}
