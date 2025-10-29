FROM golang:1.24-bookworm AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o cozystay ./cmd/api

EXPOSE 4000

CMD [ "/app/cozystay" ]