# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application (explicit output path)
RUN go build -o /app/bin/main ./main.go

# Final stage
FROM debian:stable-slim

WORKDIR /app

# Create bin directory and copy binary
RUN mkdir -p /app/bin
COPY --from=builder /app/bin/main /app/bin/
COPY gqlgen.yml ./

# Expose the application port
EXPOSE 8080

# Command to run the application (using absolute path)
CMD ["/app/bin/main"]