#!/usr/bin/env bash
# Launches unit tests
GO111MODULE=on go test -v github.com/ouroboros-crypto/node/x/structure/keeper
GO111MODULE=on go test -v github.com/ouroboros-crypto/node/x/paramining/keeper