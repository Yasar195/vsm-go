# Multi-stage build for optimized Lambda container

# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -tags lambda.norpc \
    -ldflags="-w -s" \
    -o main .

# Runtime stage
FROM public.ecr.aws/lambda/go:1.22

# Copy the binary from builder stage
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/

# Copy any additional files if needed (like .env for local testing)
# COPY .env ${LAMBDA_TASK_ROOT}/

# Set the CMD to your handler
CMD [ "main" ]