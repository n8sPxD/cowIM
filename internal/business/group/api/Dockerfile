FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/business/group/api

RUN go build -o group group.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/business/group/api/

COPY --from=builder $SERVICE_DIR/configs/group.yaml                 /app/etc/group.yaml
COPY --from=builder $SERVICE_DIR/group  /app/group

WORKDIR /app

EXPOSE 20005

ENTRYPOINT ["/app/group"]
