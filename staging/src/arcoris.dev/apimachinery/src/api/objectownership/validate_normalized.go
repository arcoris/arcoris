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

package objectownership

import "reflect"

// ValidateNormalized checks that doc is valid and already canonical.
//
// ValidateNormalized preserves raw validation errors from Validate. For valid
// raw documents, it compares doc with Normalize(doc) and reports
// ErrNotNormalized when sorting, merging, deduplication, pruning, version
// canonicalization, or nil-vs-empty canonicalization would change the document.
func ValidateNormalized(doc Document) error {
	if err := Validate(doc); err != nil {
		return err
	}

	normalized, err := Normalize(doc)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(doc, normalized) {
		return errorAt(
			pathDocument,
			ErrNotNormalized,
			ErrorReasonNotNormalized,
			"object ownership document is not normalized",
		)
	}

	return nil
}
