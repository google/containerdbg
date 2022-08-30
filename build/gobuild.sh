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

VERBOSE=${VERBOSE:-"0"}
V=""
if [[ "${VERBOSE}" == "1" ]];then
    V="-x"
    set -x
fi

SCRIPTPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

OUT=${1:?"output path"}
shift

set -e

BUILD_GOOS=${GOOS:-linux}
BUILD_GOARCH=${GOARCH:-amd64}
GOBINARY=${GOBINARY:-go}
GOPKG="$GOPATH/pkg"
BUILDINFO=${SCRIPTPATH}/.buildinfo
STATIC=${STATIC:-1}
LDFLAGS=${LDFLAGS:--extldflags -static}
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
# Split GOBUILDFLAGS by spaces into an array called GOBUILDFLAGS_ARRAY.
IFS=' ' read -r -a GOBUILDFLAGS_ARRAY <<< "$GOBUILDFLAGS"

GCFLAGS=${GCFLAGS:-}
export CGO_ENABLED=${CGO_ENABLED:-0}


if [[ "${STATIC}" !=  "1" ]];then
    LDFLAGS=""
fi

cat ${BUILDINFO}

# BUILD LD_EXTRAFLAGS
LD_EXTRAFLAGS=""

while read -r line; do
    LD_EXTRAFLAGS="${LD_EXTRAFLAGS} -X ${line}"
done < "${BUILDINFO}"

OPTIMIZATION_FLAGS="-trimpath"
if [ "${DEBUG}" == "1" ]; then
    OPTIMIZATION_FLAGS=""
fi

BUILDLOG=`mktemp`
time GOOS=${BUILD_GOOS} GOARCH=${BUILD_GOARCH} ${GOBINARY} build \
      ${V} "${GOBUILDFLAGS_ARRAY[@]}" ${GCFLAGS:+-gcflags "${GCFLAGS}"} \
      ${GOTAGS:+-tags "${GOTAGS}"} \
      -o "${OUT}" \
      ${OPTIMIZATION_FLAGS} \
      -pkgdir="${GOPKG}/${BUILD_GOOS}_${BUILD_GOARCH}" \
      -ldflags "${LDFLAGS} ${LD_EXTRAFLAGS}" "${@}" 2>&1 >$BUILDLOG
BUILDRESULT=$?

cat $BUILDLOG
exit $BUILDRESULT
