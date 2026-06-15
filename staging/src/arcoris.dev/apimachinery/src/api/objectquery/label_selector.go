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

// LabelSelector is a canonical AND-set of label requirements.
//
// Empty selectors match every item. Non-empty selectors require every
// requirement to match the item's metadata labels.
type LabelSelector struct {
	// requirements are sorted by key, operator, and values.
	requirements []LabelRequirement
}

// IsZero reports whether s has no label requirements.
func (s LabelSelector) IsZero() bool {
	return len(s.requirements) == 0
}

// Requirements returns the selector requirements in canonical order.
//
// The returned slice and every requirement inside it are detached from the
// selector. Mutating the slice or values returned from each requirement does
// not affect s.
func (s LabelSelector) Requirements() []LabelRequirement {
	return cloneLabelRequirements(s.requirements)
}

// clone returns a detached selector.
func (s LabelSelector) clone() LabelSelector {
	return LabelSelector{requirements: cloneLabelRequirements(s.requirements)}
}

// cloneLabelRequirements returns detached label requirements.
func cloneLabelRequirements(requirements []LabelRequirement) []LabelRequirement {
	if len(requirements) == 0 {
		return nil
	}

	out := make([]LabelRequirement, len(requirements))
	for i, req := range requirements {
		out[i] = req.clone()
	}

	return out
}
