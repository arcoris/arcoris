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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Validate checks v against descriptor starting at the semantic root path.
//
// The descriptor is expected to have been validated at construction,
// registration, or catalog boundaries. Validate does not call types.ValidateResolved
// for every payload; it performs local defensive descriptor checks while
// traversing the value.
func Validate(v value.Value, descriptor types.Descriptor, opts Options) error {
	return New(opts).Validate(v, descriptor)
}

// ValidateAt checks v against descriptor starting at path.
//
// The supplied path is preserved in diagnostics. This lets standalone payload
// validation start at "$" while object/surface validation can start at a
// semantic base such as "$.desired" or "$.observed".
//
// ValidateAt shares Validate's descriptor-preparation boundary: descriptors are
// expected to be prevalidated, and this function only performs local defensive
// checks needed during concrete value traversal.
func ValidateAt(path fieldpath.Path, v value.Value, descriptor types.Descriptor, opts Options) error {
	return New(opts).ValidateAt(path, v, descriptor)
}
