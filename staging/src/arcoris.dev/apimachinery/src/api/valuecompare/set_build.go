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

import "arcoris.dev/apimachinery/api/fieldpath"

// setAt builds the one-path set used for modified or opaque unknown fields.
func setAt(path fieldpath.Path) (fieldpath.Set, error) {
	set, err := fieldpath.NewSet(path)

	if err != nil {
		return fieldpath.Set{}, wrapAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"field path cannot be stored in a set",
			err,
		)
	}

	return set, nil
}
