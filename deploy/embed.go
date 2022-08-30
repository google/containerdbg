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

package deploy

import (
	"embed"
	"io/ioutil"
	"os"
	"path"
)

//go:embed *
var Deployment embed.FS

func expandToDir(rootdir string, toDir string) error {

	if err := os.MkdirAll(toDir, os.FileMode(0700)); err != nil {
		return err
	}

	entries, err := Deployment.ReadDir(rootdir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := expandToDir(path.Join(rootdir, entry.Name()), path.Join(toDir, entry.Name())); err != nil {
				return err
			}
		} else {
			data, err := Deployment.ReadFile(path.Join(rootdir, entry.Name()))
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(path.Join(toDir, entry.Name()), data, os.FileMode(0700)); err != nil {
				return err
			}

		}
	}

	return nil
}

func ExpandToDir(toDir string) error {
	return expandToDir(".", toDir)
}
