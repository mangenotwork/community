@echo off
set GOOS=windows

rmdir /s /q .\dist\burrow

mkdir .\dist\burrow
mkdir .\dist\burrow\config

go mod tidy

go build -o .\dist\burrow\burrow.exe ..\cmd\burrow\main.go

copy ..\config\burrow.yaml .\dist\burrow\config\burrow.yaml
