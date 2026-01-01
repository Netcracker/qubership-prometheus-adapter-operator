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
FROM registry.k8s.io/prometheus-adapter/prometheus-adapter:v0.12.0 AS source

FROM alpine:3.23.2

WORKDIR /
COPY --from=source /adapter ./
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
