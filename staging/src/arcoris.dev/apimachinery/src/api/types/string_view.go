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

// StringView exposes read-only DescriptorString payload data.
type StringView struct {
	// payload is a detached copy of the string descriptor payload.
	payload stringPayload
}

// AsString returns a string view when desc is DescriptorString.
func (desc Descriptor) AsString() (StringView, bool) {
	if desc.code != DescriptorString {
		return StringView{}, false
	}

	return StringView{payload: cloneStringPayload(desc.string)}, true
}

// MinBytes returns the string minimum UTF-8 byte length rule.
func (v StringView) MinBytes() (int, bool) {
	return v.payload.minBytes.value, v.payload.minBytes.set
}

// MaxBytes returns the string maximum UTF-8 byte length rule.
func (v StringView) MaxBytes() (int, bool) {
	return v.payload.maxBytes.value, v.payload.maxBytes.set
}

// MinRunes returns the string minimum rune count rule.
func (v StringView) MinRunes() (int, bool) {
	return v.payload.minRunes.value, v.payload.minRunes.set
}

// MaxRunes returns the string maximum rune count rule.
func (v StringView) MaxRunes() (int, bool) {
	return v.payload.maxRunes.value, v.payload.maxRunes.set
}

// Pattern returns the string pattern rule.
func (v StringView) Pattern() (string, bool) {
	return v.payload.pattern, v.payload.hasPattern
}

// Enum returns detached string enum values.
func (v StringView) Enum() []string {
	return cloneSlice(v.payload.enum)
}
