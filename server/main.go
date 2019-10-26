package main

import (
	servicenow "github.com/mattermost/mattermost-plugin-servicenow/server/plugin"
	"github.com/mattermost/mattermost-server/plugin"
)

func main() {
	plugin.ClientMain(&servicenow.Plugin{})
}
