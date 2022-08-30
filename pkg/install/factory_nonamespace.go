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

package install

import (
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type noNamespaceClientConfig struct {
	inner clientcmd.ClientConfig
}

type noNamespaceFactory struct {
	cmdutil.Factory
}

func NewNoNamespaceFactory(inner cmdutil.Factory) cmdutil.Factory {
	return &noNamespaceFactory{
		Factory: inner,
	}
}

func (f *noNamespaceFactory) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return &noNamespaceClientConfig{
		inner: f.Factory.ToRawKubeConfigLoader(),
	}
}
func (config *noNamespaceClientConfig) RawConfig() (clientcmdapi.Config, error) {
	return config.RawConfig()
}

func (config *noNamespaceClientConfig) ClientConfig() (*restclient.Config, error) {
	return config.ClientConfig()
}

func (config *noNamespaceClientConfig) Namespace() (string, bool, error) {
	ns, _, err := config.inner.Namespace()

	return ns, false, err

}

// ConfigAccess returns the rules for loading/persisting the config.
func (config *noNamespaceClientConfig) ConfigAccess() clientcmd.ConfigAccess {
	return config.ConfigAccess()
}
