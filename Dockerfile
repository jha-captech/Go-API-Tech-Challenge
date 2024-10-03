# Use the official Golang image to create a build environment
FROM golang:1.23.0-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app and output to /app/myapp
RUN go build -o /app/myapp ./cmd/api

# Use a minimal image to run the binary
FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/myapp .

# Expose port 8000
EXPOSE 8000

# Run the built Go binary
CMD ["./myapp"]
