package main

import (
	"fmt"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu-community/sensu-plugin-sdk/templates"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"os"
	"strings"
)

type HandlerConfig struct {
	sensu.PluginConfig
	rocketchatUrl       			string
	rocketchatChannel   			string
	rocketchatUsername  			string
	rocketchatPassword				string
	rocketchatDescriptionTemplate 	string
}

const (
	url					= "url"
	channel             = "channel"
	username            = "username"
	descriptionTemplate = "description-template"

	defaultChannel  = "general"
	defaultUsername = "servicep"
	defaultTemplate = "{{ .Check.Output }}"
)

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-rocketchat-handler",
			Short:    "The Sensu Go Rocketchat handler for notifying a channel",
			Keyspace: "abes140377/plugins/rocketchat/config",
		},
	}

	rocketchatConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      url,
			Env:       "ROCKETCHAT_URL",
			Argument:  url,
			Shorthand: "w",
			Usage:     "The webhook url to send messages to",
			Value:     &config.rocketchatUrl,
		},
		{
			Path:      channel,
			Env:       "ROCKETCHAT_CHANNEL",
			Argument:  channel,
			Shorthand: "c",
			Default:   defaultChannel,
			Usage:     "The channel to post messages to",
			Value:     &config.rocketchatChannel,
		},
		{
			Path:      username,
			Env:       "ROCKETCHAT_USERNAME",
			Argument:  username,
			Shorthand: "u",
			Default:   defaultUsername,
			Usage:     "The username that messages will be sent as",
			Value:     &config.rocketchatUsername,
		},
		{
			Path:      descriptionTemplate,
			Env:       "ROCKETCHAT_DESCRIPTION_TEMPLATE",
			Argument:  descriptionTemplate,
			Shorthand: "t",
			Default:   defaultTemplate,
			Usage:     "The Rocketchat notification output template, in Golang text/template format",
			Value:     &config.rocketchatDescriptionTemplate,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, rocketchatConfigOptions, checkArgs, sendMessage)
	goHandler.Execute()
}

func checkArgs(_ *corev2.Event) error {
	// Support deprecated environment variables
	if url := os.Getenv("ROCKETCHAT_URL"); url != "" {
		config.rocketchatUrl = url
	}
	if channel := os.Getenv("ROCKETCHAT_CHANNEL"); channel != "" && config.rocketchatChannel == defaultChannel {
		config.rocketchatChannel = channel
	}
	if username := os.Getenv("ROCKETCHAT_USERNAME"); username != "" && config.rocketchatUsername == defaultUsername {
		config.rocketchatUsername = username
	}

	if len(config.rocketchatUrl) == 0 {
		return fmt.Errorf("--%s or ROCKETCHAT_URL environment variable is required", url)
	}

	return nil
}

func formattedEventAction(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "RESOLVED"
	default:
		return "ALERT"
	}
}

func eventKey(event *corev2.Event) string {
	return fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name)
}

func eventSummary(event *corev2.Event, maxLength int) string {
	output := chomp(event.Check.Output)
	if len(event.Check.Output) > maxLength {
		output = output[0:maxLength] + "..."
	}
	return fmt.Sprintf("%s:%s", eventKey(event), output)
}

func chomp(s string) string {
	return strings.Trim(strings.Trim(strings.Trim(s, "\n"), "\r"), "\r\n")
}

func formattedMessage(event *corev2.Event) string {
	return fmt.Sprintf("%s - %s", formattedEventAction(event), eventSummary(event, 100))
}

func messageColor(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "good"
	case 2:
		return "danger"
	default:
		return "warning"
	}
}

func messageStatus(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "Resolved"
	case 2:
		return "Critical"
	default:
		return "Warning"
	}
}

func sendMessage(event *corev2.Event) error {
	client := rest.Client{Protocol: "http", Host: "chat.dzbw.de", Port: "80"}
	credentials := &models.UserCredentials{Name: "servicep", Email: "servicep@dzbw.de", Password: "servicep"}

	client.Login(credentials)

	general := &models.Channel{ID: "GENERAL", Name: "general"}

	description, errEvalTemplate := templates.EvalTemplate("description", config.rocketchatDescriptionTemplate, event)
	if errEvalTemplate != nil {
		fmt.Errorf("Error processing template: %s", errEvalTemplate)
	}

	errSend := client.Send(general, description)

	return errSend
}