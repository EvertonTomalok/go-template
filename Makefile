ENV_VARS =$(shell grep -v "^\#" .env | xargs)

swag:
	swag init

server:
	go run . server

build:
	go build -o exe main.go

envvars:
	@echo "export $(ENV_VARS)"