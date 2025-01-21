# Use the official Go image as a base
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the entire application code, including subdirectories
COPY . .

# Build the Go application
RUN go build -o main ./cmd/api

# Command to run the executable
CMD ["./main"] 