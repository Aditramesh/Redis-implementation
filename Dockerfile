# Build stage
FROM golang:1.26.2 AS builder

WORKDIR /app

# Download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app .

# Runtime stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 9999

CMD ["./app"]