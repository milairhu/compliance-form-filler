FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

COPY . .
RUN go mod download

# Build the CLI binary
RUN go build -o compliance-form-filler ./cmd/form-filler/main.go

FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/compliance-form-filler .

# Set the default command
ENTRYPOINT ["./compliance-form-filler", "answer", "--source-file", "${SOURCE_FILE}" , "--output-file", "${OUTPUT_FILE}", "--qdrant-url", "${QDRANT_URL}", "--llm-url", "${LLM_URL}", "--llm-model", "${LLM_MODEL}", "--embedding-api-url", "${EMBEDDING_API_URL}"]
