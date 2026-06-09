// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectlifecycle

import (
	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectvalidation"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// config stores constructor options before validation.
type config struct {
	// store is the required commit dependency.
	store objectstore.Store

	// resources is the required descriptor lookup dependency.
	resources ResourceResolver

	// resolver is optional unless descriptors contain DescriptorRef values.
	resolver types.Resolver

	// desiredValidator is required for every objectvalidation.Plan.
	desiredValidator objectvalidation.SurfaceValidator[value.Value]

	// observedValidator is optional until a request carries Observed data.
	observedValidator objectvalidation.SurfaceValidator[value.Value]

	// applyOptions supplies objectapply traversal knobs.
	applyOptions objectapply.Options
}

// newConfig applies and validates constructor options.
func newConfig(opts []Option) (config, error) {
	var cfg config
	for _, opt := range opts {
		if opt == nil {
			return config{}, errorFor(0, ErrorReasonInvalidExecutor, objectstore.Key{}, ErrInvalidExecutor, ErrNilOption)
		}
		opt(&cfg)
	}

	if err := validateConfig(cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}
