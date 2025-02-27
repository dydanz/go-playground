# Build stage
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app
ENV DB_HOST=host.docker.internal
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=go_cursor
ENV REDIS_HOST=host.docker.internal
ENV REDIS_PORT=6379

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the entire application code, including subdirectories
COPY . .

# Build the Go application
RUN go build -o loyalty_engine ./cmd/api

# Final stage
FROM alpine:3.21
RUN apk add --no-cache libc6-compat

# Set the working directory
WORKDIR /app

# Copy migrations folder and binary
COPY --from=builder /app/server/migrations ./server/migrations
COPY --from=builder /app/server/docs ./server/docs
COPY --from=builder /app/loyalty_engine .
COPY --from=builder /app/.env .

# Expose the port your app runs on
EXPOSE 8080

# Run the application
CMD ["./loyalty_engine"]
