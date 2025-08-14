FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags lambda.norpc -ldflags="-w -s" -o main .

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/
CMD [ "main" ]
