FROM golang:1.16-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o slo-computer

# Use a minimal alpine image for the final image
FROM alpine:3.14

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/slo-computer /app/slo-computer

# Create a non-root user to run the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["/app/slo-computer"] 