.PHONY: build-alpine up

build-alpine:
	GOOS=linux GOARCH=amd64 go build -o server *.go

up: build-alpine
	docker rm -f kafka-sender | true && \
	docker-compose up -d --build --force-recreate; rm server