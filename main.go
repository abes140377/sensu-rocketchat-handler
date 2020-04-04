package main

import (
	"fmt"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu-community/sensu-plugin-sdk/templates"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	neturl "net/url"
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
	password            = "password"
	descriptionTemplate = "description-template"

	defaultUrl  = "https://open.rocket.chat/"
	defaultChannel  = "sandbox"
	defaultUsername = "sensu"
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
			Default:   defaultUrl,
			Usage:     "The Rocketchat Server URL to send messages to",
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
			Path:      password,
			Env:       "ROCKETCHAT_PASSWORD",
			Argument:  password,
			Shorthand: "p",
			Usage:     "The password of the user",
			Value:     &config.rocketchatPassword,
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
	if password := os.Getenv("ROCKETCHAT_PASSWORD"); password != "" {
		config.rocketchatPassword = password
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
	u, parseErr := neturl.Parse(config.rocketchatUrl)

	if parseErr != nil {
		fmt.Errorf("Error parsing url : %s", config.rocketchatUrl)
	}

	client := rest.Client{Protocol: u.Scheme, Host: u.Host, Port: u.Port()}
	// credentials := &models.UserCredentials{Name: config.rocketchatUsername, Email: "servicep@dzbw.de", Password: "servicep"}
	credentials := &models.UserCredentials{Name: config.rocketchatUsername, Password: config.rocketchatPassword}

	loginErr := client.Login(credentials)

	if loginErr != nil {
		fmt.Errorf("Error login with username : %s", config.rocketchatUsername)
	}

	// channel := &models.Channel{ID: "GENERAL", Name: "channel"}
	channel := &models.Channel{Name: config.rocketchatChannel}

	description, errEvalTemplate := templates.EvalTemplate("description", config.rocketchatDescriptionTemplate, event)
	if errEvalTemplate != nil {
		fmt.Errorf("Error processing template: %s", errEvalTemplate)
	}

	errSend := client.Send(channel, description)

	return errSend
}