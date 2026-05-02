FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/laundryServer
RUN apk add --no-cache sqlite

FROM alpine:latest
RUN apk add --no-cache sqlite-libs
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/auxiliary ./auxiliary
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./server", "-config", "config.yaml"]