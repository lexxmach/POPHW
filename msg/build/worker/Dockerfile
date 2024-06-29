# syntax=docker/dockerfile:1

FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o worker cmd/worker/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /build/worker ./
COPY ./build/worker/config.json ./config.json

CMD ./worker -c config.json