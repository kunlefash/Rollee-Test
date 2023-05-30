FROM golang:latest

WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port that the server will listen on
EXPOSE 8545

# Run the Go application
CMD ["./main"]
