#!/bin/bash

rm -rf ./dist/burrow
mkdir -p ./dist/burrow
mkdir -p ./dist/burrow/config

go mod tidy

go build -o ./dist/burrow/burrow ./cmd/burrow/main.go

copy ./config/burrow.yaml ./dist/burrow/config/burrow.yaml