# syntax=docker/dockerfile:1

FROM golang:1.22-alpine as builder 

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o client cmd/client/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /build/client ./
COPY ./build/client/config.json ./config.json

CMD ./client -c config.json