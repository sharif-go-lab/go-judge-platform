# --- Builder Stage ---
FROM golang:latest AS builder

# Set working directory
WORKDIR /code

# Cache and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go server binary (static, optimized)
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags "-s -w" -o server ./cmd/server

# --- Final Stage ---
FROM alpine:latest

# Install certificates (for HTTPS, external APIs, etc.)
RUN apk add --no-cache ca-certificates

# Create app directory
WORKDIR /app

# Copy binary and application assets
COPY --from=builder /code/server .
COPY --from=builder /code/config.yaml .
COPY --from=builder /code/templates ./templates
COPY --from=builder /code/static ./static
COPY --from=builder /code/migrations ./migrations

# Expose the server port
EXPOSE 8080

# Launch the server
ENTRYPOINT ["./server"]