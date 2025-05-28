FROM golang:1.24.2-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api-server main.go 

# Production stage
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api-server .

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 40000

CMD ["./api-server"]