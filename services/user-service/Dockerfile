FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./services/user-service/cmd/api

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/services/user-service/migrations ./migrations

EXPOSE 8081 9091

CMD ["./main"] 