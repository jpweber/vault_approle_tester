#!/bin/sh
env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o vault_approle_test_linux main.go