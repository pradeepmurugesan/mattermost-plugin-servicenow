package connect

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {

	expected := "[click here](http://localhost:8080/plugins/mattermost-plugin-servicenow/oauth/connect) to connect with ServiceNow"
	siteURL := "http://localhost:8080"
	config := &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	}

	assert.Equal(t, Execute(config), expected)

}
