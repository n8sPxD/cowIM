FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/apis/infos/api

RUN go build -o infos infos.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/apis/infos/api/

COPY --from=builder $SERVICE_DIR/etc/infos.yaml /app/etc/infos.yaml
COPY --from=builder $SERVICE_DIR/infos          /app/infos

WORKDIR /app

EXPOSE 20001

ENTRYPOINT ["/app/infos"]
