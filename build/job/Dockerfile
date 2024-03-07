ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

################################################################################

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build

RUN --mount=type=cache,target=/var/cache/apk \
  apk --update add --no-cache ca-certificates tzdata \
  && update-ca-certificates

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 \
  go build -v \
  -ldflags "-s -w" \
  -buildvcs=false -trimpath -buildmode=exe \
  -o /bin/job ./cmd/job

ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  appuser

################################################################################

FROM alpine:${ALPINE_VERSION} as compress

ARG UPX_SOURCE_AMD64=https://github.com/upx/upx/releases/download/v4.2.2/upx-4.2.2-amd64_linux.tar.xz
ARG UPX_ARCHIVE_AMD64=upx-4.2.2-amd64_linux.tar.xz
ARG UPX_DIRECTORY_AMD64=upx-4.2.2-amd64_linux

ARG UPX_SOURCE_ARM64=https://github.com/upx/upx/releases/download/v4.2.2/upx-4.2.2-arm64_linux.tar.xz
ARG UPX_ARCHIVE_ARM64=upx-4.2.2-arm64_linux.tar.xz
ARG UPX_DIRECTORY_ARM64=upx-4.2.2-arm64_linux

RUN if [ "$(uname -m)" = "x86_64" ]; then \
    wget ${UPX_SOURCE_AMD64} \
    && tar -xf ${UPX_ARCHIVE_AMD64} \
    && cd ${UPX_DIRECTORY_AMD64}; \
  else \
    wget ${UPX_SOURCE_ARM64} \
    && tar -xf ${UPX_ARCHIVE_ARM64} \
    && cd ${UPX_DIRECTORY_ARM64}; \
  fi \
    && chmod +x ./upx \
    && mv ./upx /bin/upx

RUN mkdir -p /compress

COPY --from=build --chmod=770 /bin/job /compress/job-original

WORKDIR /compress

RUN upx --best --lzma job-original -o job-compressed

################################################################################

FROM scratch as original

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build --chmod=770 /bin/job /service/job

USER appuser

ENTRYPOINT ["./service/job"]

################################################################################

FROM scratch as reduced

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=compress --chmod=770 /compress/job-compressed /service/job

USER appuser

ENTRYPOINT ["./service/job"]
