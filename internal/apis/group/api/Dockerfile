FROM n8spxd/cow-im:base AS builder

WORKDIR /build/internal/apis/group/api

RUN go build -o group group.go

FROM golang:alpine

ARG SERVICE_DIR=/build/internal/apis/group/api/

COPY --from=builder $SERVICE_DIR/etc/group.yaml /app/etc/group.yaml
COPY --from=builder $SERVICE_DIR/group          /app/group

WORKDIR /app

EXPOSE 20001

ENTRYPOINT ["/app/group"]
