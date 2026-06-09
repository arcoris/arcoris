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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Extract returns the semantic paths explicitly mentioned by v from "$".
//
// The descriptor is expected to have been validated at construction,
// registration, or catalog boundaries. Extract does not call types.ValidateResolved
// on every payload; it performs local defensive checks needed for read-only
// traversal and stable path construction.
func Extract(
	val value.Value,
	descriptor types.Descriptor,
	opts Options,
) (fieldpath.Set, error) {
	return ExtractAt(fieldpath.RootPath(), val, descriptor, opts)
}

// ExtractAt returns the semantic paths explicitly mentioned by v from path.
//
// The supplied path is preserved in every extracted child path. This lets
// standalone payload extraction start at "$" while future object/surface layers
// can start at semantic bases such as "$.desired" or "$.observed".
func ExtractAt(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	opts Options,
) (fieldpath.Set, error) {
	if err := path.Validate(); err != nil {
		return fieldpath.Set{}, wrapAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"base field path is invalid",
			err,
		)
	}

	run := newExtractor(opts)
	return run.extract(path, val, descriptor, 0)
}
