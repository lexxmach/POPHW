# syntax=docker/dockerfile:1

FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/server/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /build/server ./
COPY ./build/server/config.json ./config.json

CMD ./server -c config.json

EXPOSE 8080