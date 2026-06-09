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
	"unicode/utf8"
)

// validateString checks DescriptorString length, pattern, and enum descriptor rules.
//
// Pattern text is compiled only for descriptor validation. The descriptor keeps
// the original pattern string so future codecs and schema exporters can choose
// their own representation without depending on Go regexp internals.
func validateString(desc Descriptor, path string) error {
	if err := validateLengthLimits(desc.string.minBytes, desc.string.maxBytes, path+".bytes"); err != nil {
		return err
	}
	if err := validateLengthLimits(desc.string.minRunes, desc.string.maxRunes, path+".runes"); err != nil {
		return err
	}

	if desc.string.hasPattern {
		compiled, err := regexp.Compile(desc.string.pattern)

		if err != nil {
			return descriptorErrorf(
				path+".pattern",
				ErrInvalidDescriptor,
				DescriptorErrorReasonInvalidPattern,
				"string pattern %q is not a valid regexp: %v",
				desc.string.pattern,
				err,
			)
		}

		for i, value := range desc.string.enum {
			if !compiled.MatchString(value) {
				return descriptorErrorf(
					fmt.Sprintf("%s.enum[%d]", path, i),
					ErrInvalidDescriptor,
					DescriptorErrorReasonEnumPatternMismatch,
					"string enum value %q does not match pattern %q",
					value,
					desc.string.pattern,
				)
			}
		}
	}

	if index, value, ok := firstDuplicate(desc.string.enum); ok {
		return descriptorErrorf(
			path+".enum",
			ErrInvalidDescriptor,
			DescriptorErrorReasonDuplicateEnum,
			"string enum values must be unique; duplicate value %q at index %d",
			value,
			index,
		)
	}

	for i, value := range desc.string.enum {
		byteLen := len(value)
		runeLen := utf8.RuneCountInString(value)

		if desc.string.minBytes.set && byteLen < desc.string.minBytes.value {
			return descriptorErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidDescriptor,
				DescriptorErrorReasonEnumBelowMinimum,
				"string enum value %q byte length %d is below minimum byte length %d",
				value,
				byteLen,
				desc.string.minBytes.value,
			)
		}

		if desc.string.maxBytes.set && byteLen > desc.string.maxBytes.value {
			return descriptorErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidDescriptor,
				DescriptorErrorReasonEnumAboveMaximum,
				"string enum value %q byte length %d is above maximum byte length %d",
				value,
				byteLen,
				desc.string.maxBytes.value,
			)
		}

		if desc.string.minRunes.set && runeLen < desc.string.minRunes.value {
			return descriptorErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidDescriptor,
				DescriptorErrorReasonEnumBelowMinimum,
				"string enum value %q rune count %d is below minimum rune count %d",
				value,
				runeLen,
				desc.string.minRunes.value,
			)
		}

		if desc.string.maxRunes.set && runeLen > desc.string.maxRunes.value {
			return descriptorErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidDescriptor,
				DescriptorErrorReasonEnumAboveMaximum,
				"string enum value %q rune count %d is above maximum rune count %d",
				value,
				runeLen,
				desc.string.maxRunes.value,
			)
		}
	}

	return nil
}
