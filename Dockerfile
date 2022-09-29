ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /builder
WORKDIR /builder

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./server ./cmd/main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /
COPY --from=builder /builder/server .
COPY --from=builder /builder/common.yaml .
COPY --from=builder /builder/.env .

EXPOSE 9000
