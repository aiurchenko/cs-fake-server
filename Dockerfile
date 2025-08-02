FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/server .
COPY .env .env
EXPOSE 27015/udp
ENTRYPOINT ["./server"]
