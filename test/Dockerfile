# Stage 1: Build the Go binary
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY test/go.mod test/go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY test/ .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create a lightweight Docker image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port on which the microservice will listen
EXPOSE 8080

# Run the Go binary
CMD ["./main"]
