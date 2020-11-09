#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=1.0 -extldflags '-static -s'" -o bin/dev.exe main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=1.0 -extldflags '-static -s'" -o bin/dev main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=1.0 -extldflags '-static -s'" -o bin/mac/dev main.go
