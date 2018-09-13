.PHONY: build-alpine up

build-alpine:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

up: build-alpine
	docker rm -f kafka-sender | true && \
	docker-compose up -d --build --force-recreate; rm server