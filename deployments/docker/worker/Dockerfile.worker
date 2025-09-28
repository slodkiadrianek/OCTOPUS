FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o worker ./cmd/worker

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/worker .

COPY .env .

CMD ["./worker"]