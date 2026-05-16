.PHONY: swagger run build

swagger:
	swag init -g main.go -o docs --parseDependency --parseInternal --dir ./cmd/server,./internal/handlers,./internal/apidocs

build:
	go build -o bin/server.exe ./cmd/server

run:
	go run ./cmd/server
