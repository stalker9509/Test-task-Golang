FROM golang:1.22 AS builder
WORKDIR /app

COPY . .

RUN go build -o app ./cmd/app

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]