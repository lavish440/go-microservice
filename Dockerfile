FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates wget

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o calc-server ./server/main.go

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.25 && \
    wget -qO /bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe


FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/calc-server /app/calc-server
COPY --from=builder /bin/grpc_health_probe /app/grpc_health_probe
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

EXPOSE 50051

USER nonroot:nonroot

HEALTHCHECK --interval=10s --timeout=2s --retries=3 \
  CMD ["/app/grpc_health_probe", "-addr=:50051"]

ENTRYPOINT ["/app/calc-server"]
