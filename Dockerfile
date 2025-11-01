FROM golang:1.24-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o cozystay ./cmd/api

FROM scratch AS production

WORKDIR /prod

COPY --from=builder /build/cozystay ./

EXPOSE 4000

CMD [ "/prod/cozystay" ]
