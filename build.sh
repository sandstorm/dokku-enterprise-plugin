#!/usr/bin/env bash

export GOOS=linux
export GOARCH=386


#go build -o bin-build/commands gitlab.sandstorm.de/infrastructure/dokku-plugin/commands/
#go build -o bin-build/user-auth gitlab.sandstorm.de/infrastructure/dokku-plugin/user-auth/
go build -o bin-build/nginx-pre-reload github.com/sandstorm/dokku-enterprise-plugin/plugins/nginx-pre-reload/
go build -o bin-build/post-deploy github.com/sandstorm/dokku-enterprise-plugin/plugins/post-deploy/
go build -o bin-build/commands github.com/sandstorm/dokku-enterprise-plugin/plugins/commands/
