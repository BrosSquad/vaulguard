FROM golang:1.15-alpine as builder
RUN apk add --no-cache git sqlite make gcc curl g++ sqlite-dev

RUN curl -sL https://taskfile.dev/install.sh | sh

WORKDIR /go/src/github.com/BrosSquad/vaulguard
COPY . .


RUN task test
RUN task build



FROM alpine:3.12

COPY --from=builder /go/src/github.com/BrosSquad/vaulguard/bin/vaulguard /vaulguard/vaulguard

EXPOSE 8000 

ENTRYPOINT ["/vaulguard/vaulguard"]
CMD ["-config","/etc/vaulguard/config.yml", "-port", "8000"] 
