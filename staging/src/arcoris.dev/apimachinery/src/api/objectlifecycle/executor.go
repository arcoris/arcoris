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
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// ResourceResolver resolves already registered resource definitions.
//
// It is intentionally the same method set as resource.Resolver. The local name
// documents objectlifecycle ownership without depending on a concrete
// resourcecatalog.Catalog implementation.
type ResourceResolver interface {
	resource.Resolver
}

// Executor coordinates lifecycle operations over one object store.
//
// Executor has no background goroutines, no package-level registration, and no
// mutable global state. Its dependencies are configured explicitly by NewExecutor.
type Executor struct {
	// store commits already-computed live object state.
	store objectstore.Store

	// resources resolves resource contracts by GVK or GVR.
	resources ResourceResolver

	// resolver resolves structural type references for validation and apply.
	resolver types.Resolver

	// desiredValidator validates the required Desired surface.
	desiredValidator objectvalidation.SurfaceValidator[value.Value]

	// observedValidator validates Observed when a version defines it and input carries it.
	observedValidator objectvalidation.SurfaceValidator[value.Value]

	// applyOptions carries resolver/depth knobs into objectapply.
	applyOptions objectapply.Options
}

// NewExecutor constructs a lifecycle executor from explicit dependencies.
func NewExecutor(opts ...Option) (*Executor, error) {
	cfg, err := newConfig(opts)
	if err != nil {
		return nil, err
	}

	return &Executor{
		store:             cfg.store,
		resources:         cfg.resources,
		resolver:          cfg.resolver,
		desiredValidator:  cfg.desiredValidator,
		observedValidator: cfg.observedValidator,
		applyOptions:      cfg.applyOptions,
	}, nil
}
