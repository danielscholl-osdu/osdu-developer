FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mise-mock .

# Create a minimal image
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/mise-mock .

EXPOSE 9002

ENTRYPOINT ["./mise-mock"]