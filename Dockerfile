FROM golang:1.25.4 AS build
WORKDIR /go/src/app
COPY . .

RUN set -eux; \
    CGO_ENABLED=0 GO111MODULE=on go install .; \
    go run github.com/google/go-licenses@latest save ./... --save_path=/notices; \
    mkdir -p /mounts/config;

FROM ghcr.io/greboid/dockerbase/nonroot:1.20250803.0
COPY --from=build /go/bin/tsp /tsp
COPY --from=build /notices /notices
COPY --from=build --chown=65532:65532 /mounts /
VOLUME /config
ENTRYPOINT ["/tsp", "--tailscale-config-dir=/config"]