FROM golang:1.24.4-alpine
WORKDIR /app
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w" -o api-server ./cmd/server

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/api-server .

EXPOSE 3000
ENTRYPOINT ["/app/api-server"]