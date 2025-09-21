# Go image
FROM golang:1.24.3

# Set working directory
WORKDIR /app

# Copy go files
COPY . .

# Download dependencies
RUN go mod download

# Build
RUN go build -o main .

# Run the binary
CMD ["./main"]
