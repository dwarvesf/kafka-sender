FROM golang:1.10
RUN mkdir -p /go/src/github.com/dwarvesf/kafka-sender
WORKDIR /go/src/github.com/dwarvesf/kafka-sender
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server main.go

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/dwarvesf/kafka-sender/server server
EXPOSE 3030
CMD ["/server"]