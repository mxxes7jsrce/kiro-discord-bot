# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bot .

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
# Set timezone to Asia/Taipei for local use
ENV TZ=Asia/Taipei
WORKDIR /app
COPY --from=builder /app/bot .
VOLUME ["/data", "/projects"]
CMD ["./bot"]
