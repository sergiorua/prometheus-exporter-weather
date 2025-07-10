FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-exporter ./cmd/weather-exporter

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/weather-exporter .

# Copy config file
COPY --from=builder /app/configs/config.yaml ./configs/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./weather-exporter"]