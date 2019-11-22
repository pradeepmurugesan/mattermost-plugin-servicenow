package plugin

import (
	"github.com/mattermost/mattermost-plugin-servicenow/server/handlers/api/hello"
	"github.com/mattermost/mattermost-plugin-servicenow/server/handlers/api/incident"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/plugin"
	"net/http"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/hello":
		hello.SayHello(w, r)
	case "/incident":
		if err := incident.PublishIncident(w, r, models.PluginContext{API: p.API, BotUserID: p.BotUserID}); err != nil {
			mlog.Error(err.Message)
		}

	default:
		http.NotFound(w, r)
	}
}
