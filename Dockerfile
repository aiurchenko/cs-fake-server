FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 27015/udp
ENTRYPOINT ["./server"]
