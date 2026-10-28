# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build

WORKDIR /src
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY rogue/ ./rogue/
RUN go test ./... && go build -trimpath -ldflags="-s -w" -o /out/rogue ./rogue/cmd
RUN mkdir /data && chown 10001:10001 /data

FROM scratch

ENV TERM=xterm-256color
COPY --from=build /out/rogue /rogue
COPY --from=build --chown=10001:10001 /data /data
USER 10001:10001
WORKDIR /data

ENTRYPOINT ["/rogue", "-save-dir", "/data"]
