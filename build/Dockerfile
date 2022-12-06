FROM golang:1.17-alpine AS builder
WORKDIR /build/
COPY ./ ./
RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux go build -o gigi

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/gigi ./
CMD ["./gigi"]
