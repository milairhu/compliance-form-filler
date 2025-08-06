FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o compliance-form-filler ./cmd/form-filler/main.go

FROM alpine:latest

# Install curl and jq
RUN apk add --no-cache curl jq

COPY --from=builder /app/compliance-form-filler .
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
