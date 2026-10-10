FROM golang:1.21-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY main.go .

RUN go mod init algonode
RUN go mod tidy

RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Внутренний порт остается 8080 (наружу в Windows он проброшен как 8085)
EXPOSE 8080
CMD ["./main"]