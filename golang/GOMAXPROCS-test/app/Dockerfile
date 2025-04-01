# Build stage
FROM golang:1.24 AS builder

# Set CGO_ENABLED to 1 for CGO support
ENV CGO_ENABLED=1

WORKDIR /app

# Copy and build the application
COPY . .
RUN go build -o app .

# Final stage
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/app /app

# Set the entrypoint
ENTRYPOINT ["/app"]