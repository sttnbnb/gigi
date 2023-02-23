FROM golang:1.20-alpine AS builder
WORKDIR /build/
COPY ./ ./
RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux go build -o gigi

FROM scratch
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/gigi ./
CMD ["./gigi"]
