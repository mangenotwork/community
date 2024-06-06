@echo off
set GOOS=linux
set GOARCH=amd64

rmdir /s /q .\dist\burrow_linux

mkdir .\dist\burrow_linux
mkdir .\dist\burrow_linux\config

go mod tidy

go build -o .\dist\burrow_linux\burrow ..\cmd\burrow\main.go

copy ..\config\burrow.yaml .\dist\burrow_linux\config\burrow.yaml
