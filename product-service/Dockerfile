FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/.env /app

# Use ARG for build-time variables
ARG SERVER_PORT=8081
ENV SERVER_PORT=${SERVER_PORT}

# Use the ARG value in EXPOSE
EXPOSE ${SERVER_PORT}

CMD ["./main"]