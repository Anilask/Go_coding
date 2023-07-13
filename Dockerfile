# Start from the official Go image as the base
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /go

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the rest of the source code to the working directory
COPY . .

# Build the Go application
RUN go Build

# Set the container port to listen on
EXPOSE 8000

# Set the command to run when the container starts
CMD ["./go-api"]
