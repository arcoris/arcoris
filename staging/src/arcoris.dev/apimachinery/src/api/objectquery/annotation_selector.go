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

// AnnotationSelector is a canonical AND-set of annotation requirements.
//
// Empty selectors match every item. Non-empty selectors require every
// requirement to match the item's metadata annotations.
type AnnotationSelector struct {
	// requirements are sorted by key, operator, and values.
	requirements []AnnotationRequirement
}

// IsZero reports whether s has no annotation requirements.
func (s AnnotationSelector) IsZero() bool {
	return len(s.requirements) == 0
}
