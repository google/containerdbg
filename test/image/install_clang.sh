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

CLANG_VERSION=llvmorg-13.0.1

apt-get update
apt-get -y install cmake

git clone --depth 1 -b "${CLANG_VERSION}" https://github.com/llvm/llvm-project.git

cd llvm-project

mkdir build

cd build

cmake -DCMAKE_BUILD_TYPE=Release -DLLVM_ENABLE_PROJECTS=clang -G "Unix Makefiles" ../llvm

make -j"$(nproc)"

make install

cd ../..

rm -rf llvm-project

ln -s /usr/local/bin/llvm-strip /usr/local/bin/llvm-strip-13
