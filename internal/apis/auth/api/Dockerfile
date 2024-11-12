FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/apis/auth/api

RUN go build -o auth auth.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/apis/auth/api/

COPY --from=builder $SERVICE_DIR/etc/auth.yaml /app/etc/auth.yaml
COPY --from=builder $SERVICE_DIR/auth          /app/auth

WORKDIR /app

EXPOSE 20001

ENTRYPOINT ["/app/auth"]