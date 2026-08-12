# Copyright 2025 NetCracker Technology Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# hadolint global ignore=DL3018
FROM --platform=$BUILDPLATFORM golang:1.26.5-alpine3.24@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
# renovate: datasource=github-releases depName=kubernetes-sigs/prometheus-adapter versioning=semver
ARG PROMETHEUS_ADAPTER_VERSION=v0.12.0

WORKDIR /workspace

RUN apk add --no-cache \
        git \
    && git clone --depth 1 --branch "${PROMETHEUS_ADAPTER_VERSION}" \
        https://github.com/kubernetes-sigs/prometheus-adapter.git .

RUN go mod edit \
        -require=go.opentelemetry.io/otel@v1.44.0 \
        -require=go.opentelemetry.io/otel/metric@v1.43.0 \
        -require=go.opentelemetry.io/otel/sdk@v1.43.0 \
        -require=go.opentelemetry.io/otel/trace@v1.43.0 \
        -require=golang.org/x/crypto@v0.52.0 \
        -require=golang.org/x/net@v0.57.0 \
        -require=golang.org/x/oauth2@v0.34.0 \
        -require=golang.org/x/text@v0.40.0 \
        -require=google.golang.org/grpc@v1.82.1 \
    && go mod tidy

RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" GO111MODULE=on \
    go build -mod=readonly -trimpath -o /adapter sigs.k8s.io/prometheus-adapter/cmd/adapter

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

WORKDIR /
COPY --from=builder /adapter ./
COPY entrypoint.sh ./

ENV USER_UID=1001 \
    USER_NAME=adapter \
    GROUP_NAME=adapter

RUN apk add --no-cache \
        bash \
        procps \
        inotify-tools \
        tini \
    # Cleanup after install utilities
    && rm -rf /var/cache/apk/* \
    # Add user
    && addgroup ${GROUP_NAME} && adduser -D -G ${GROUP_NAME} -u ${USER_UID} ${USER_NAME} \
    # Grant permissions to run entrypoint
    && chmod +x entrypoint.sh

USER ${USER_UID}

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
