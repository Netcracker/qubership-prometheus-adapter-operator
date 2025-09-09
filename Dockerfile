# Copyright 2025 NetCracker Technology Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.25.1-alpine3.22 AS builder
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy the Go sources
COPY api/ api/
COPY controllers/ controllers/
COPY main.go main.go
COPY go.* /workspace/

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go work sync

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GO111MODULE=on go build -a -o manager main.go

# Use alpine tiny images as a base
FROM alpine:3.22.1

ENV USER_UID=1001 \
    USER_NAME=operator

WORKDIR /
COPY --from=builder --chown=${USER_UID} /workspace/manager .

USER ${USER_UID}

ENTRYPOINT ["/manager"]
