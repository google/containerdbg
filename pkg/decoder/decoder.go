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

package decoder

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	goyaml "gopkg.in/yaml.v3"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/k8s"
)

type Options struct {
	DefaultGVK  *schema.GroupVersionKind
	MutateFuncs []MutateFunc
}

type DecodeOption func(*Options)

type HandlerFunc func(ctx context.Context, obj k8s.Object) error

type MutateFunc func(obj k8s.Object) error

func DefaultGVK(def *schema.GroupVersionKind) DecodeOption {
	return func(o *Options) {
		o.DefaultGVK = def
	}
}

func WithMutation(mut MutateFunc) DecodeOption {
	return func(o *Options) {
		o.MutateFuncs = append(o.MutateFuncs, mut)
	}
}

func DecodeAny(manifest io.Reader, options ...DecodeOption) (k8s.Object, error) {
	opts := Options{}
	for _, opt := range options {
		opt(&opts)
	}
	k8sDecoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer().Decode
	b, err := io.ReadAll(manifest)
	if err != nil {
		return nil, err
	}

	runtimeObj, _, err := k8sDecoder(b, opts.DefaultGVK, nil)
	if runtime.IsNotRegisteredError(err) {
		runtimeObj = &unstructured.Unstructured{}
		if err := yaml.Unmarshal(b, runtimeObj); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	obj, ok := runtimeObj.(k8s.Object)
	if !ok {
		return nil, fmt.Errorf("couldn't convert %T to Object", obj)
	}
	for _, patch := range opts.MutateFuncs {
		if err := patch(obj); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func DecodeEach(ctx context.Context, manifest io.Reader, handlerFn HandlerFunc, options ...DecodeOption) error {
	decoder := yaml.NewYAMLReader(bufio.NewReader(manifest))
	for {
		b, err := decoder.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		obj, err := DecodeAny(bytes.NewReader(b), options...)
		if err != nil {
			return err
		}
		if err := handlerFn(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func DecodeUnstructured(unstructured *unstructured.Unstructured) (k8s.Object, error) {
	data, _ := goyaml.Marshal(unstructured.Object)
	k8sDecoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer().Decode
	runtimeObj, _, err := k8sDecoder(data, nil, nil)
	if runtime.IsNotRegisteredError(err) {
		return unstructured, nil
	} else if err != nil {
		return nil, err
	}

	obj, ok := runtimeObj.(k8s.Object)
	if !ok {
		return nil, fmt.Errorf("couldn't convert %T to Object", obj)
	}

	return obj, nil
}
