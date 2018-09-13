FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY server /server
EXPOSE 3030
ENTRYPOINT ["/server"]