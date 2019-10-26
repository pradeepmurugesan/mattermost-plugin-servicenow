package plugin

import (
	"github.com/mattermost/mattermost-plugin-servicenow/server/handlers/api/hello"
	"github.com/mattermost/mattermost-server/plugin"
	"net/http"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/hello":
		hello.SayHello(w, r)
	default:
		http.NotFound(w, r)
	}
}
