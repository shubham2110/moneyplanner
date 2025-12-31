# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with CGO disabled for scratch compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o moneyplanner .

# Final stage
FROM scratch

# Copy the binary from builder
COPY --from=builder /app/moneyplanner /moneyplanner

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["/moneyplanner"]