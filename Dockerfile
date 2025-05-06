FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Initialize module and get dependencies
RUN go mod init github.com/yourusername/todo-api
RUN go get github.com/gorilla/mux
RUN go get github.com/joho/godotenv
RUN go get github.com/lib/pq
RUN go mod tidy

# Build the application
RUN go build -o main ./cmd/api/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]