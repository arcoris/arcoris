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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Compare reports semantic changes between two present payload values at "$".
//
// The descriptor is expected to have been validated before comparison. Compare
// performs local traversal checks only: invalid zero values, invalid descriptor
// views, kind mismatches, unresolved TypeRef values, and invalid ListMap keys.
func Compare(
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	opts Options,
) (Result, error) {
	return CompareAt(fieldpath.RootPath(), oldValue, newValue, descriptor, opts)
}

// CompareAt reports semantic changes between two present payload values at path.
//
// The supplied base path is preserved in every returned path and diagnostic.
// This lets callers compare nested surfaces without rewriting root-based
// results.
func CompareAt(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	opts Options,
) (Result, error) {
	if err := path.Validate(); err != nil {
		return Result{}, wrapAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"base field path is invalid",
			err,
		)
	}

	run := newComparer(opts)
	return run.compare(
		path,
		presentOperand(oldValue),
		presentOperand(newValue),
		descriptor,
		0,
	)
}
