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

WORKDIR /labstore

COPY server/go.mod ./server/go.mod
RUN cd server/ && go mod download

COPY client/go.mod ./client/go.mod
RUN cd client/ && go mod download

COPY cli/go.mod ./cli/go.mod
RUN cd cli/ && go mod download

COPY go.work go.work
RUN go work sync

COPY . .
COPY --from=node_builder /web/build ./server/pkg/router/assets


RUN GOOS=$TARGETOS GOARCH=$TARGETARCH cd cli/ && go build -o ../bin/labstore ./cmd/labstore

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
