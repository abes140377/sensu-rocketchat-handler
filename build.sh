#!/usr/bin/env bash
#

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -extldflags "-static"' -o sensu-rocketchat-handler