FROM golang:1.24.1-alpine AS builder

# Set the working directory
WORKDIR /app

COPY . .

# Build the CLI binary
RUN go build -o compliance-form-filler ./cmd/form-filler/main.go

FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/compliance-form-filler .
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x entrypoint.sh

# Set the default command
ENTRYPOINT ["./entrypoint.sh"]
