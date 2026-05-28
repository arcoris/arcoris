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

import "regexp"

// validateString checks TypeString length, pattern, and enum descriptor rules.
//
// Pattern text is compiled only for descriptor validation. The descriptor keeps
// the original pattern string so future codecs and schema exporters can choose
// their own representation without depending on Go regexp internals.
func validateString(t Type, path string) error {
	if err := validateLengthLimits(t.string.minLen, t.string.maxLen, path+".len"); err != nil {
		return err
	}
	if t.string.hasPattern {
		compiled, err := regexp.Compile(t.string.pattern)
		if err != nil {
			return typeError(path+".pattern", ErrInvalidType)
		}
		for _, value := range t.string.enum {
			if !compiled.MatchString(value) {
				return typeError(path+".enum", ErrInvalidType)
			}
		}
	}
	if hasDuplicates(t.string.enum) {
		return typeError(path+".enum", ErrInvalidType)
	}
	for _, value := range t.string.enum {
		if t.string.minLen.set && len(value) < t.string.minLen.value {
			return typeError(path+".enum", ErrInvalidType)
		}
		if t.string.maxLen.set && len(value) > t.string.maxLen.value {
			return typeError(path+".enum", ErrInvalidType)
		}
	}
	return nil
}
