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

from kubernetes import client, config

config.load_incluster_config() #download config inside of cluster

apiservice_client = client.ApiregistrationV1Api()
api_service_name = 'v1beta1.custom.metrics.k8s.io'
apiservice = apiservice_client.read_api_service(api_service_name)

available_condition = None
for condition in apiservice.status.conditions:
    if condition.type == 'Available':
        available_condition = condition
        break

if available_condition:
    print(
        f"API Service: {api_service_name}, Status: {available_condition.status}, "
        f"Message: {available_condition.message}, Reason: {available_condition.reason}"
    )
    if available_condition.status == 'False':
        exit(1)
    else:
        exit(0)
else:
    print(f"No availability condition found for API Service: {api_service_name}")
    exit(1)