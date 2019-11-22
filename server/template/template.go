package template

import (
	"bytes"
	"fmt"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"text/template"
)

const incident string = `
### New incident created

Created By : {{.SysCreatedBy}}
Impact : {{.Impact}}
Priority : {{.Priority}}
Description: {{.ShortDescription}}
`

var templateMap = make(map[string]string)

func init() {
	templateMap["incident"] = incident
}

//RenderTemplate renders the given template if present in the templateMap
func RenderTemplate(name string, data interface{}) (string, error) {
	temp, ok := templateMap[name]

	if !ok {
		return "", &models.Error{Message: fmt.Sprintf("Couldn't find the template with name %s", name)}
	}
	var output bytes.Buffer
	t, err := template.New(name).Parse(temp)
	if err != nil {
		return "", err
	}

	err = t.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
