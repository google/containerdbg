// Copyright 2021 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package imagehelpers

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func GetImageEntryPoint(imageName string) ([]string, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, err
	}

	image, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		var errLocal error
		image, errLocal = daemon.Image(ref)
		if errLocal != nil {
			return nil, fmt.Errorf("didn't find the image %s in local docker: %s, or remote: %s", imageName, errLocal, err)
		}
	}

	config, err := image.ConfigFile()
	if err != nil {
		return nil, err
	}

	return append(config.Config.Entrypoint, config.Config.Cmd...), nil
}
