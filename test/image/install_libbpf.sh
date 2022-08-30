#!/bin/bash -e
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

LIBBPF_VERSION=0.7.0

apt-get -y install libelf-dev

mkdir libbpf
pushd libbpf

wget https://github.com/libbpf/libbpf/archive/refs/tags/v${LIBBPF_VERSION}.tar.gz
tar xf v${LIBBPF_VERSION}.tar.gz
cd libbpf-${LIBBPF_VERSION}/src
make
make install
popd
rm -rf libbpf
