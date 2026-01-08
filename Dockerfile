# --- Build Stage ---
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

# 1. Install git (not always needed but good for dependencies)
RUN apk --no-cache add git

# Pre-copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build Monolithic App
RUN go build -o unihub cmd/server/main.go

# --- Run Stage ---
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary
COPY --from=builder /app/unihub .

# Copy configs
COPY --from=builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Start command
CMD ["./unihub"]
