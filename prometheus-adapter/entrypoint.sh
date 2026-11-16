#!/bin/bash

# Copyright 2024-2025 NetCracker Technology Corporation
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

# Existed prometheus-adapter ConfigMap name
MOUNTED_PATH="/etc/adapter"

# Runs prometheus-adapter
manage_adapter() {
    while true; do
        echo "Running prometheus-adapter ..."
        "$@" & # Run adapter as a daemon
        adapter_pid=$!
        echo "Ran prometheus-adapter with adapter_pid $adapter_pid"
        wait $adapter_pid

        # Kill adapter
        if kill -0 $adapter_pid >&/dev/null; then
            echo "Kill job_pid $adapter_pid"
            kill $adapter_pid
        fi
        
        echo "Exiting ..."
        break
    done
}

main() {
    # Run prometheus-adapter
    manage_adapter "$@"
}

main "$@"
