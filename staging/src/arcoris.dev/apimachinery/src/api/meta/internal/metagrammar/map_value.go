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

package metagrammar

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// MapValueOptions controls metadata map value validation.
type MapValueOptions struct {
	// AllowEmpty controls whether an empty value is accepted.
	AllowEmpty bool
	// MaxLength limits the value in bytes when greater than zero.
	MaxLength int
	// Strict restricts values to the small label-value character set.
	Strict bool
}

// ValidateMapValue validates a metadata map value.
func ValidateMapValue(value string, opts MapValueOptions) *Violation {
	if value == "" && !opts.AllowEmpty {
		return violation(ReasonEmptyValue, "map value must be non-empty")
	}

	if opts.MaxLength > 0 && len(value) > opts.MaxLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("map value length must be <= %d bytes", opts.MaxLength),
		)
	}

	if HasControl(value) {
		return violation(ReasonInvalidCharacter, "map value must not contain control bytes")
	}

	if !opts.Strict {
		return nil
	}

	if value != "" && !lexical.IsASCIIAlnum(value[0]) {
		return violation(
			ReasonInvalidEdge,
			"map value must start with an ASCII letter or digit",
		)
	}

	if value != "" && !lexical.IsASCIIAlnum(value[len(value)-1]) {
		return violation(
			ReasonInvalidEdge,
			"map value must end with an ASCII letter or digit",
		)
	}

	for i := 0; i < len(value); i++ {
		b := value[i]
		if isStrictMapValueByte(b) {
			continue
		}

		return violation(
			ReasonInvalidCharacter,
			fmt.Sprintf("map value contains invalid byte %q at index %d", b, i),
		)
	}

	return nil
}

// isStrictMapValueByte reports whether b is allowed inside strict label values.
func isStrictMapValueByte(b byte) bool {
	return lexical.IsASCIIAlnum(b) || b == '-' || b == '_' || b == '.'
}
