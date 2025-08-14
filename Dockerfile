FROM golang:1.24-alpine AS builder
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build for x86_64 Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags lambda.norpc -ldflags="-w -s" -o bootstrap .

# Lambda runtime
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/bootstrap ${LAMBDA_TASK_ROOT}/
# CMD is optional — Lambda will run bootstrap by default
CMD ["bootstrap"]
