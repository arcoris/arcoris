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

package valuevalidation

import (
	"regexp"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateString checks TypeString kind, length, pattern, and enum constraints.
func (v *validator) validateString(path fieldpath.Path, val value.Value, descriptor types.Type) {
	if !v.requireKind(path, val, value.KindString, descriptor.Code()) {
		return
	}

	text, _ := val.String()
	stringView, ok := descriptor.String()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not a string",
		)
		return
	}

	v.validateStringLength(path, text, stringView)
	v.validateStringPattern(path, text, stringView)
	v.validateStringEnum(path, text, stringView)
}

// validateStringLength checks byte-length rules using the same semantics as api/types.
func (v *validator) validateStringLength(path fieldpath.Path, text string, stringView types.StringView) {
	length := len(text)
	if minLength, ok := stringView.MinLen(); ok && length < minLength {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooShort,
			"string length %d is below minimum %d",
			length,
			minLength,
		)
	}

	if maxLength, ok := stringView.MaxLen(); ok && length > maxLength {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooLong,
			"string length %d is above maximum %d",
			length,
			maxLength,
		)
	}
}

// validateStringPattern checks a descriptor regexp against the concrete string.
func (v *validator) validateStringPattern(path fieldpath.Path, text string, stringView types.StringView) {
	pattern, ok := stringView.Pattern()
	if !ok {
		return
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		v.wrap(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"string descriptor pattern is invalid",
			err,
		)
		return
	}

	if !compiled.MatchString(text) {
		v.addf(
			path,
			ErrPatternMismatch,
			ErrorReasonPatternMismatch,
			"string value does not match pattern %q",
			pattern,
		)
	}
}

// validateStringEnum checks string enum membership when the descriptor has one.
func (v *validator) validateStringEnum(path fieldpath.Path, text string, stringView types.StringView) {
	enum := stringView.Enum()
	if len(enum) == 0 || containsString(enum, text) {
		return
	}

	v.add(
		path,
		ErrEnumMismatch,
		ErrorReasonEnumMismatch,
		"string value is not in enum",
	)
}

// containsString reports whether values contains target.
func containsString(values []string, target string) bool {
	for _, candidate := range values {
		if candidate == target {
			return true
		}
	}

	return false
}
