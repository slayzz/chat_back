FROM golang:1.14.2-alpine3.11

RUN apk add --update tzdata \
    bash wget curl git

RUN mkdir -p /opt/todoapp
WORKDIR /opt/todoapp

ENV GO111MODULE=on

COPY go.mod .
RUN go mod download
