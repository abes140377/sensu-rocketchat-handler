[![Bonsai Asset Badge](https://img.shields.io/badge/Sensu%20Slack%20Handler-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/sensu/sensu-rocketchat-handler)
![goreleaser](https://github.com/abes140377/sensu-rocketchat-handler/workflows/goreleaser/badge.svg)

# Sensu Rocketchat Handler

- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
  - [Check definition](#check-definition)

## Overview


The [Sensu Slack Handler][0] is a [Sensu Event Handler][3] that sends event data
to a configured Slack channel.

## Usage examples

Help:

```
Usage:
  sensu-rocketchat-handler [flags]
  sensu-rocketchat-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -w, --rocketchat-url string   The url of the Rocket.Chat server to send messages to (default "https://open.rocket.chat/")
  -c, --channel string          The channel to post messages to (default "sandbox")
  -i, --icon-url string         A URL to an image to use as the user avatar (default "https://www.sensu.io/img/sensu-logo.png")
  -u, --username string         The username that messages will be sent as (default "sensu")
  -t, --descriptionTemplate     The Rocketchat notification output template, in Golang text/template format
  -h, --help                    help for sensu-rocketchat-handler
```

## Configuration

### Asset registration

Assets are the best way to make use of this handler. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 or later, you can use the following command to add the asset:

`sensuctl asset add abes140377/sensu-rocketchat-handler`

If you're using an earlier version of sensuctl, you can download the asset definition from [this project's
Bonsai Asset Index page][6].

### Handler definition

Create the handler using the following handler definition:

```yml
---
api_version: core/v2
type: Handler
metadata:
  namespace: default
  name: rocketchat
spec:
  type: pipe
  command: sensu-rocketchat-handler -channel 'sandbox' --username 'sensu'
  filters:
  - is_incident
  runtime_assets:
  - sensu/sensu-rocketchat-handler
  secrets:
  - name: ROCKETCHAT_PASSWORD
    secret: rocketchat-password
  timeout: 10
```