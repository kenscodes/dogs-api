# Build stage
FROM golang:1.21-bookworm AS builder

WORKDIR /app

# Install gcc for CGO (required for SQLite)
RUN apt-get update && apt-get install -y gcc

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy backend source
COPY backend/ ./backend/

# Copy frontend files
COPY frontend/ ./frontend/

# Build the application
RUN cd backend && CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o dogs-api main.go

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/backend/dogs-api .

# Copy frontend
COPY --from=builder /app/frontend ./frontend

# Copy dogs.json
COPY --from=builder /app/backend/dogs.json ./dogs.json

# Create data directory
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Run the application with persistent database
ENV DB_PATH=/data/dogs.db
ENV FRONTEND_PATH=/root/frontend
CMD ["./dogs-api"]
