# SPDX-License-Identifier: Apache-2.0
FROM golang:1.24-alpine-3.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/magefile/mage@latest \
  && go install golang.org/x/tools/cmd/goimports@latest \
  && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

COPY . .
RUN mage build

FROM alpine:3.24
LABEL org.opencontainers.image.source=https://github.com/srediag/srediag
RUN addgroup -S app && adduser -S -G app app
WORKDIR /home/app
COPY --from=builder /app/bin/srediag .
COPY --from=builder /app/configs ./configs
USER app
ENTRYPOINT ["./srediag"]
CMD ["--config", "configs/config.yaml"]
