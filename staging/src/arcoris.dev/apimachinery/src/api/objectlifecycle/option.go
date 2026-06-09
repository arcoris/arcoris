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

// Option configures NewExecutor.
type Option func(*config)

// WithStore sets the objectstore commit dependency.
func WithStore(store objectstore.Store) Option {
	return func(cfg *config) {
		cfg.store = store
	}
}

// WithResourceResolver sets the resource contract resolver.
func WithResourceResolver(resources ResourceResolver) Option {
	return func(cfg *config) {
		cfg.resources = resources
	}
}

// WithTypeResolver sets the structural type reference resolver.
func WithTypeResolver(resolver types.Resolver) Option {
	return func(cfg *config) {
		cfg.resolver = resolver
	}
}

// WithDesiredValidator sets the required Desired surface validator.
func WithDesiredValidator(validator objectvalidation.SurfaceValidator[value.Value]) Option {
	return func(cfg *config) {
		cfg.desiredValidator = validator
	}
}

// WithObservedValidator sets the optional Observed surface validator.
func WithObservedValidator(validator objectvalidation.SurfaceValidator[value.Value]) Option {
	return func(cfg *config) {
		cfg.observedValidator = validator
	}
}

// WithApplyOptions sets objectapply options used by Apply.
//
// ApplyRequest.Force is per-request state and replaces the Force field from
// these options for each Apply call. Resolver and MaxDepth are preserved.
func WithApplyOptions(opts objectapply.Options) Option {
	return func(cfg *config) {
		cfg.applyOptions = opts
	}
}
