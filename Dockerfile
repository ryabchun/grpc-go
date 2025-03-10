FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o grpc_service ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/grpc_service /app/
EXPOSE 50051
CMD ["./grpc_service"]
