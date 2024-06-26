FROM golang:1.22 as builder

WORKDIR /go/src

RUN go install github.com/cosmtrek/air@latest

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

ENV DOTENV_PATH=/go/src/.env

CMD ["air"]
