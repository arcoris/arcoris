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

import "arcoris.dev/apimachinery/api/fieldpath"

// setBuilder collects extracted paths and canonicalizes them once.
type setBuilder struct {
	paths []fieldpath.Path
}

// Add records one extracted path.
func (b *setBuilder) Add(path fieldpath.Path) {
	b.paths = append(b.paths, path)
}

// AddSet records every path from set without allocating a detached slice.
func (b *setBuilder) AddSet(set fieldpath.Set) {
	set.ForEach(func(_ int, path fieldpath.Path) bool {
		b.Add(path)
		return true
	})
}

// Build returns the canonical extracted field set.
func (b *setBuilder) Build(path fieldpath.Path) (fieldpath.Set, error) {
	set, err := fieldpath.NewSet(b.paths...)
	if err != nil {
		return fieldpath.Set{}, wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"extracted field path cannot be stored in a set",
			err,
		)
	}

	return set, nil
}

// setAt builds the one-path field set used for scalar, null, empty composite,
// atomic-list, and opaque-unknown-field extraction.
func setAt(path fieldpath.Path) (fieldpath.Set, error) {
	set, err := fieldpath.NewSet(path)
	if err != nil {
		return fieldpath.Set{}, wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"field path cannot be stored in a set",
			err,
		)
	}

	return set, nil
}
