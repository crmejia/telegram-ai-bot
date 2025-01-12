ARG GO_VERSION=1
#FROM golang:${GO_VERSION}-bookworm as builder
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


#FROM debian:bookworm
FROM alpine:latest

#RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
