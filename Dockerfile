FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gopot ./cmd/gopot

FROM alpine:3.21
RUN addgroup -S gopot && adduser -S gopot -G gopot
WORKDIR /app
RUN chown gopot:gopot /app
COPY --from=builder /build/gopot /usr/local/bin/gopot
USER gopot
EXPOSE 2223
ENTRYPOINT ["gopot"]
