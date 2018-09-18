.PHONY: build-alpine up

IMAGE?=hub.dwarvesf.com/yggdrasil/kafka-sender
VERSION?=latest

build-alpine:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

up: build-alpine
	docker rm -f kafka-sender | true && \
	docker-compose up -d --build --force-recreate; rm server

package: build-alpine
	docker build -t $(IMAGE):$(VERSION) .

ship: package
	docker push $(IMAGE):$(VERSION)

deploy:
	kubectl create -f kubernetes/deployment.yaml
	kubectl create -f kubernetes/service.yaml

update:
	kubectl apply -f kubernetes/deployment.yaml