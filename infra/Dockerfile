# =================
# Node builder
# =================

FROM --platform=$BUILDPLATFORM node:25-alpine3.22 AS node_builder

WORKDIR /web

COPY /web/package.json /web/package-lock.json ./
RUN npm ci

COPY web .

RUN npm run build

# =================
# Go builder
# =================

FROM --platform=$BUILDPLATFORM golang:1.25-alpine3.22 AS go_builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /labstore

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=node_builder /web/build ./server/router/assets


ARG GOARM=7
RUN set -eux; \
    if [ "$TARGETARCH" = "arm" ]; then \
        export GOARM="${TARGETVARIANT#v}"; \
    fi; \
    GO_CONSTANTS="github.com/IllumiKnowLabs/labstore/server/constants"; \
    LDFLAGS="-X ${GO_CONSTANTS}.GitTag=$(git describe --tags --abbrev=0 2>/dev/null) \
        -X ${GO_CONSTANTS}.GitCommit=$(git rev-parse --short HEAD) \
        -X ${GO_CONSTANTS}.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
        -X ${GO_CONSTANTS}.Builder=docker"; \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags embedassets -ldflags "$LDFLAGS" -o bin/labstore ./cmd/labstore

# =================
# LabStore CLI
# =================

FROM alpine:3.22

COPY --from=go_builder /labstore/bin/labstore /usr/local/bin/labstore

RUN apk add --no-cache curl

RUN addgroup -S labstore && adduser -S labstore -G labstore

RUN mkdir /data /keys
RUN chown -R labstore:labstore /data /keys
RUN chmod 0700 /keys

VOLUME /data
VOLUME /keys

USER labstore

ENTRYPOINT ["labstore"]
