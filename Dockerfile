FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
# Remove GOARCH=amd64 to let Docker handle architecture automatically
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN adduser -D appuser
WORKDIR /app
COPY --from=builder /app/main .
RUN chown appuser:appuser /app/main
USER appuser
EXPOSE 8080
CMD ["./main"]