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

package types

// stringPayload stores TypeString constraints.
//
// Pattern is stored as text rather than *regexp.Regexp so descriptors remain
// serializable by future codecs and exporters. Validation may compile the
// pattern, but the descriptor itself does not own runtime regex state.
type stringPayload struct {
	// minLen is the inclusive minimum string length.
	minLen limit[int]
	// maxLen is the inclusive maximum string length.
	maxLen limit[int]

	// pattern stores the textual regular expression for string descriptors.
	pattern string
	// hasPattern distinguishes an explicitly empty pattern from an unset rule.
	hasPattern bool

	// enum stores accepted string literals in declaration order.
	enum []string
}

// cloneStringPayload detaches string enum values.
func cloneStringPayload(p stringPayload) stringPayload {
	p.enum = cloneSlice(p.enum)

	return p
}

// emptyStringPayload reports whether p has no configured TypeString state.
func emptyStringPayload(p stringPayload) bool {
	return !p.minLen.set && !p.maxLen.set && !p.hasPattern && p.pattern == "" && len(p.enum) == 0
}
