# Use an official Go runtime as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy the entire project into the container
COPY . .
# Generate Wire code
RUN go install github.com/google/wire/cmd/wire@latest
RUN wire ./cmd/server/runner

# Build the Go application
RUN go build -o myapp ./cmd/server

# Set the command to run your application
CMD ["./myapp"]
