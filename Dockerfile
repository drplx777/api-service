FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-server .
EXPOSE 3000
ENTRYPOINT ["/app/api-server"]