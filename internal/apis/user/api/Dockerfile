FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/apis/user/api

RUN go build -o user user.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/apis/user/api/

COPY --from=builder $SERVICE_DIR/etc/user.yaml /app/etc/user.yaml
COPY --from=builder $SERVICE_DIR/user          /app/user

WORKDIR /app

EXPOSE 20001

ENTRYPOINT ["/app/user"]
