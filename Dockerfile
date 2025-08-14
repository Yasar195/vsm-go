# Use the official AWS Lambda Go base image
FROM public.ecr.aws/lambda/go:1.24

# Set the working directory
WORKDIR ${LAMBDA_TASK_ROOT}

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -tags lambda.norpc -o main .

# Copy the built binary to the Lambda runtime directory
RUN cp main ${LAMBDA_RUNTIME_DIR}

# Set the CMD to your handler (the function name in your Go code)
CMD [ "main" ]