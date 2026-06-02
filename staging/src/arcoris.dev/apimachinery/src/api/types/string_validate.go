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

import (
	"fmt"
	"regexp"
)

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
			return typeErrorf(
				path+".pattern",
				ErrInvalidType,
				TypeErrorReasonInvalidPattern,
				"string pattern %q is not a valid regexp: %v",
				t.string.pattern,
				err,
			)
		}

		for i, value := range t.string.enum {
			if !compiled.MatchString(value) {
				return typeErrorf(
					fmt.Sprintf("%s.enum[%d]", path, i),
					ErrInvalidType,
					TypeErrorReasonEnumPatternMismatch,
					"string enum value %q does not match pattern %q",
					value,
					t.string.pattern,
				)
			}
		}
	}

	if index, value, ok := firstDuplicate(t.string.enum); ok {
		return typeErrorf(
			path+".enum",
			ErrInvalidType,
			TypeErrorReasonDuplicateEnum,
			"string enum values must be unique; duplicate value %q at index %d",
			value,
			index,
		)
	}

	for i, value := range t.string.enum {
		if t.string.minLen.set && len(value) < t.string.minLen.value {
			return typeErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidType,
				TypeErrorReasonEnumBelowMinimum,
				"string enum value %q length %d is below minimum length %d",
				value,
				len(value),
				t.string.minLen.value,
			)
		}

		if t.string.maxLen.set && len(value) > t.string.maxLen.value {
			return typeErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidType,
				TypeErrorReasonEnumAboveMaximum,
				"string enum value %q length %d is above maximum length %d",
				value,
				len(value),
				t.string.maxLen.value,
			)
		}
	}

	return nil
}
