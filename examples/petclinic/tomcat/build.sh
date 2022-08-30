#!/bin/bash
# Copyright 2022 Google LLC
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

if [[ -z "${PROJECT_ID}" ]]; then
  echo "ERROR: missing env variable PROJECT_ID"
  exit 1
fi
export VERSION="${VERSION:-latest}"

echo "building gcr.io/${PROJECT_ID}/barp-tomcat-petclinic-tomcat-petclinic-5b2b62c3:${VERSION}"
gcloud builds submit --timeout 1h -t gcr.io/"${PROJECT_ID}"/barp-tomcat-petclinic-tomcat-petclinic-5b2b62c3:"${VERSION}"
