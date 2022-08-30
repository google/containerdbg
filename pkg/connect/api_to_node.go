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

package connect

import (
	"fmt"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type tunnelOptions struct {
	RESTClient   *restclient.RESTClient
	Config       *restclient.Config
	PodClient    corev1client.PodsGetter
	Address      []string
	Port         []string
	StopChannel  <-chan struct{}
	ReadyChannel chan struct{}
	*genericclioptions.IOStreams
}

type Option interface {
	applyFor(o *tunnelOptions)
}

type stopChannel struct {
	stopChannel <-chan struct{}
}

func (s stopChannel) applyFor(o *tunnelOptions) {
	o.StopChannel = s.stopChannel
}

func WithStopChannel(stopCh <-chan struct{}) Option {
	return &stopChannel{stopChannel: stopCh}
}

type readyChannel struct {
	readyChannel chan struct{}
}

func (s readyChannel) applyFor(o *tunnelOptions) {
	o.ReadyChannel = s.readyChannel
}

func WithReadyChannel(readyCh chan struct{}) Option {
	return &readyChannel{readyChannel: readyCh}
}

type ioStreamsOption struct {
	*genericclioptions.IOStreams
}

func (s ioStreamsOption) applyFor(o *tunnelOptions) {
	o.IOStreams = s.IOStreams
}

func WithIOStreams(streams *genericclioptions.IOStreams) Option {
	return &ioStreamsOption{IOStreams: streams}
}

func (o *tunnelOptions) complete() error {
	var err error
	if o.Config == nil {
		return fmt.Errorf("rest config should be set")
	}
	if o.RESTClient == nil {
		// TODO Need to avoid editing this...
		config := *o.Config
		config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
		if config.APIPath == "" {
			config.APIPath = "/api"
		}
		if config.NegotiatedSerializer == nil {
			config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
		}
		if err := rest.SetKubernetesDefaults(&config); err != nil {
			return err
		}
		o.RESTClient, err = restclient.RESTClientFor(&config)
		if err != nil {
			return err
		}
	}

	if o.IOStreams == nil {
		o.IOStreams = &genericclioptions.IOStreams{
			Out:    io.Discard,
			ErrOut: io.Discard,
		}
	}

	if o.PodClient == nil {
		clientset, err := kubernetes.NewForConfig(o.Config)
		if err != nil {
			return err
		}
		o.PodClient = clientset.CoreV1()
	}

	if o.Address == nil {
		o.Address = []string{"localhost"}
	}

	return nil
}

// StartTunnel starts tunnel to podName in namspace on the listed ports (8080:8080 etc.)
func StartTunnel(restConfig *restclient.Config, podName string, namspace string, ports []string, opts ...Option) error {
	o := tunnelOptions{
		Config: restConfig,
		Port:   ports,
	}
	for _, opt := range opts {
		opt.applyFor(&o)
	}
	if err := o.complete(); err != nil {
		return err
	}

	req := o.RESTClient.Post().
		Resource("pods").
		Namespace(namspace).
		Name(podName).
		SubResource("portforward")
	transport, upgrader, err := spdy.RoundTripperFor(o.Config)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())

	fw, err := portforward.NewOnAddresses(dialer, o.Address, o.Port, o.StopChannel, o.ReadyChannel, o.IOStreams.Out, o.IOStreams.ErrOut)
	if err != nil {
		return err
	}

	return fw.ForwardPorts()
}
