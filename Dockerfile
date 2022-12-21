# Build stage
FROM golang:1.19.4-alpine AS builder

ENV APP_NAME=eks-shared-cost-calculator

WORKDIR /go/src

COPY . .

RUN set -euxo pipefail && \
    apk add --no-cache \
        gcc \
        musl-dev \
        ca-certificates && \
    go build -o bin/${APP_NAME}

# Release stage
FROM alpine:3.16 AS release

ENV APP_NAME=eks-shared-cost-calculator

RUN set -euxo pipefail && \
    addgroup --gid 10001 -S ${APP_NAME} && \
    adduser --uid 10001 -S ${APP_NAME} -G ${APP_NAME} --gecos=""

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder --chown=10001:10001 \
     /go/src/bin/${APP_NAME} /usr/local/bin/${APP_NAME}

USER eks-shared-cost-calculator

WORKDIR /usr/local/bin

CMD ["/usr/local/bin/eks-shared-cost-calculator"]
