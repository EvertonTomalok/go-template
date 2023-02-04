ENV_VARS =$(shell grep -v "^\#" .env | xargs)

envvars:
	@echo "export $(ENV_VARS)"

swag:
	swag init

build:
	go build -o exe main.go

.PHONY: server
server: envvars swag build
	./exe server
