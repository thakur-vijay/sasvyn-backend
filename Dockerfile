FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server
RUN go build -o cleanup ./cmd/cleanup

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/cleanup .

COPY internal/database/migrations ./internal/database/migrations

EXPOSE 8080

CMD ["./server"]