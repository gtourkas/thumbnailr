#!/usr/bin/env bash

GOPATH=~/go/src/ GOOS=linux GOARCH=amd64 go build -o ./delve/dlv github.com/go-delve/delve/cmd/dlv
