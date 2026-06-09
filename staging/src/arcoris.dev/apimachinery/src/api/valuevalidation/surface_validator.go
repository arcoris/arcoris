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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// SurfaceValidator adapts api/value validation to objectvalidation's generic
// SurfaceValidator[value.Value] method shape without importing that package.
type SurfaceValidator struct {
	// Options are copied for every validation call.
	Options Options
}

// ValidateSurface validates val against descriptor.
//
// The resolver supplied by the objectvalidation plan takes precedence when it
// is non-nil. When it is nil, the resolver already configured in Options is
// preserved.
func (v SurfaceValidator) ValidateSurface(
	val value.Value,
	descriptor types.Descriptor,
	resolver types.Resolver,
) error {
	opts := v.Options
	if resolver != nil {
		opts.Resolver = resolver
	}

	return Validate(val, descriptor, opts)
}
