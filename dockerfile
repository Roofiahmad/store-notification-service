FROM golang:1.24.13-bookworm AS builder

WORKDIR /app

ENV APP_ENV=production
ENV GIN_MODE=release

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o store-notification-service .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

ENV APP_ENV=production
ENV GIN_MODE=release

WORKDIR /app
COPY --from=builder /app/store-notification-service .

USER guest

EXPOSE 3000

CMD ["./store-notification-service"]