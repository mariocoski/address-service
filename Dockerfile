# Dockerfile
FROM golang:1.20-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the source code to the workspace
COPY . .

# Build the application
RUN go build -o /bin/app ./cmd/http

FROM golang:1.20-alpine

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@v1.27.3

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /bin/app .

# Expose the port the application listens on
EXPOSE 7000

# Start the application with hot reloading
CMD ["air"]