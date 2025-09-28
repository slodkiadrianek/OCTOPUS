FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/api .

COPY .env .

EXPOSE 3040

CMD ["./api"]
