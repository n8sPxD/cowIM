FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/msg_forward/

RUN go build -o message message.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/msg_forward

COPY --from=builder $SERVICE_DIR/etc/message.yaml /app/etc/message.yaml
COPY --from=builder $SERVICE_DIR/message          /app/message

WORKDIR /app

ENTRYPOINT ["/app/message"]