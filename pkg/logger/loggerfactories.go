// Copyright 2022 Google LLC All Rights Reserved.
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

package logger

import (
	"flag"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

var logLevel int = 0

var klogFlags flag.FlagSet
var klogInitialzationSentinal sync.Once

// glog is annoyingly registering "v" for us and there is no way
// disable that. Some parts of the code use glog and this infects the entire
// binary.
// For compatibility, if we detect that v is already bound, we save the flagset
// and read the value during first log creation.
// Remove this code when we purge glog from our codebase.
var glogCompatFlagSet *flag.FlagSet

func ensureKlogInitialized() {
	klogInitialzationSentinal.Do(func() {
		klog.InitFlags(&klogFlags)
		logLevelStr := fmt.Sprint(logLevel)
		if glogCompatFlagSet != nil {
			logLevelStr = glogCompatFlagSet.Lookup("v").Value.String()
		}
		klogFlags.Set("v", logLevelStr)
		klogFlags.Set("logtostderr", "true")
	})
}

func NewHeadlessLogger() logr.Logger {
	ensureKlogInitialized()
	return klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog))
}

func BindFlags(flags *flag.FlagSet) {
	if flags == nil {
		flags = flag.CommandLine
	}

	if flags.Lookup("v") == nil {
		flags.IntVar(&logLevel, "v", 0, "number for the log level verbosity")
	} else {
		glogCompatFlagSet = flags
	}
}
