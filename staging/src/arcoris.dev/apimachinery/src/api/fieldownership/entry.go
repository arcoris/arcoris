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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// Entry stores the fields owned by one Owner.
//
// Fields are private so constructors and State normalization can preserve owner
// validation, deterministic ordering, and immutable-by-convention boundaries.
type Entry struct {
	owner  Owner
	fields fieldpath.Set
}

// Owner returns the ownership identity for e.
func (e Entry) Owner() Owner {
	return e.owner
}

// Fields returns the semantic fields owned by e.
//
// fieldpath.Set is immutable by convention and returns detached path slices from
// its own accessors, so returning it by value preserves this package's boundary.
func (e Entry) Fields() fieldpath.Set {
	return e.fields
}

// IsEmpty reports whether e owns no paths.
func (e Entry) IsEmpty() bool {
	return e.fields.IsEmpty()
}

// Validate checks whether e can be stored in State.
func (e Entry) Validate() error {
	if err := e.owner.Validate(); err != nil {
		return wrapAt(
			"",
			ErrInvalidEntry,
			ErrorReasonInvalidEntry,
			"entry owner is invalid",
			err,
		)
	}
	if err := validateFields(e.fields, "entry field path is invalid"); err != nil {
		return err
	}

	return nil
}
