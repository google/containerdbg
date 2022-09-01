#!/bin/bash
# Copyright 2022 Google LLC All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

git_rev=${GIT_REV:-$(git rev-parse --short HEAD)}

# used by hack/gobuild.sh
echo "main.vendorVersion=${GIT_REV}"
echo "velostrata-internal.googlesource.com/containerdbg.git/pkg/build.ImageRepo=${TARGET_REPO}"
echo "velostrata-internal.googlesource.com/containerdbg.git/pkg/build.PullPolicy=${IMAGE_PULL_POLICY}"
echo "velostrata-internal.googlesource.com/containerdbg.git/pkg/build.Version=${BUILD_VERSION:-$git_rev}"
echo "velostrata-internal.googlesource.com/containerdbg.git/pkg/build.GitSha=${BUILD_VERSION:-$git_rev}"
