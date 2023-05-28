FROM golang:1.17 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o word-service .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/word-service .

EXPOSE 8080

CMD ["./word-service"]
