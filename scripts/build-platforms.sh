#!/usr/bin/env bash
# Builds everything for both nix and windows platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./build/linux.64bit/ouroborosd ./cmd/ouroborosd
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./build/linux.64bit/ouroboroscli ./cmd/ouroboroscli

GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o ./build/windows.64bit/ouroborosd.exe ./cmd/ouroborosd
GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o ./build/windows.64bit/ouroboroscli.exe ./cmd/ouroboroscli
