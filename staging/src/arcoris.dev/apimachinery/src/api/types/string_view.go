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

// StringView exposes read-only TypeString payload data.
type StringView struct {
	// payload is a detached copy of the string descriptor payload.
	payload stringPayload
}

// String returns a string view when t is TypeString.
func (t Type) String() (StringView, bool) {
	return StringView{payload: cloneStringPayload(t.string)}, t.code == TypeString
}

// MinLen returns the string minimum length rule.
func (v StringView) MinLen() (int, bool) {
	return v.payload.minLen.value, v.payload.minLen.set
}

// MaxLen returns the string maximum length rule.
func (v StringView) MaxLen() (int, bool) {
	return v.payload.maxLen.value, v.payload.maxLen.set
}

// Pattern returns the string pattern rule.
func (v StringView) Pattern() (string, bool) {
	return v.payload.pattern, v.payload.hasPattern
}

// Enum returns detached string enum values.
func (v StringView) Enum() []string {
	return cloneStrings(v.payload.enum)
}
