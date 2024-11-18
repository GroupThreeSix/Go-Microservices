#Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

# Product
FROM alpine

WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /app/.env /app

ARG GRPC_PORT=50051
ENV GRPC_PORT=${GRPC_PORT}

EXPOSE ${GRPC_PORT}

CMD ["./main"]