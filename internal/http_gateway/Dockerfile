FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/http_gateway/

RUN go build -o gateway gateway.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/http_gateway

COPY --from=builder $SERVICE_DIR/etc/gateway.yaml /app/etc/gateway.yaml
COPY --from=builder $SERVICE_DIR/gateway          /app/gateway

WORKDIR /app

EXPOSE 20001

ENTRYPOINT ["/app/gateway"]
