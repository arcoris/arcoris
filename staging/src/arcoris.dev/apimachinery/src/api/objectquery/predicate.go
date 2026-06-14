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

package objectquery

// Predicate is a validated, canonical object list item predicate.
type Predicate struct {
	// identity is the canonical identity predicate.
	identity IdentitySelector

	// labels is the canonical label predicate.
	labels LabelSelector

	// annotations is the canonical annotation predicate.
	annotations AnnotationSelector
}

// IsZero reports whether p matches every item.
func (p Predicate) IsZero() bool {
	return p.identity.IsZero() && p.labels.IsZero() && p.annotations.IsZero()
}
