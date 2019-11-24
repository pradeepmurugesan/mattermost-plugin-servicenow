package plugin

import (
	"fmt"
	"github.com/mattermost/mattermost-plugin-servicenow/server/handlers/command/hello"
	"github.com/mattermost/mattermost-plugin-servicenow/server/handlers/command/oauth/connect"
	"github.com/mattermost/mattermost-plugin-servicenow/server/handlers/command/stream"
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"strings"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "now",
		DisplayName:      "Service Now",
		Description:      "Integration with Service Now.",
		AutoComplete:     true,
		AutoCompleteDesc: "hello",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

//ExecuteCommand executes the command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""
	if len(split) > 1 {
		action = split[1]
	}

	if command != "/now" {
		return &model.CommandResponse{}, nil
	}

	switch action {
	case "hello":
		result := hello.Execute()
		p.postCommandResponse(args, fmt.Sprintf(result))
		return &model.CommandResponse{}, nil

	case "stream":
		result, err := stream.Execute(&models.PluginContext{API: p.API}, args)
		if err != nil {
			mlog.Error(err.Message)
			p.postErrorMessage(args, err)
			return nil, &model.AppError{Message: err.Message}
		}
		p.postCommandResponse(args, fmt.Sprintf(result))
		return &model.CommandResponse{}, nil

	case "connect":
		response := connect.Execute(p.API.GetConfig())
		p.postCommandResponse(args, response)
		return &model.CommandResponse{}, nil
	}

	p.postCommandResponse(args, fmt.Sprintf("Unknown action %v", action))

	return &model.CommandResponse{}, nil
}

func (p *Plugin) postErrorMessage(args *model.CommandArgs, err *models.Error) {
	p.postCommandResponse(args, err.Message)
}
