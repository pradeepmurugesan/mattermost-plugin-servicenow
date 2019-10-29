package models

import (
	"github.com/mattermost/mattermost-server/plugin"
)

// PluginContext context of the plugin
type PluginContext struct {
	API       plugin.API
	BotUserId string
}

func (p *PluginContext) SetApi(api plugin.API) {
	p.API = api
}
