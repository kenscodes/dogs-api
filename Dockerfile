# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy backend source
COPY backend/ ./backend/

# Build the application
RUN cd backend && CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o dogs-api main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

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
CMD ["./dogs-api"]
