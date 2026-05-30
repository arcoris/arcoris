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

import "fmt"

// OpaqueTokenOptions controls validation for opaque metadata scalar tokens.
type OpaqueTokenOptions struct {
	// AllowEmpty controls whether an absent token is valid at this boundary.
	AllowEmpty bool
	// MaxLength limits token size in bytes when greater than zero.
	MaxLength int
}

// ValidateOpaqueToken validates an opaque scalar token without parsing internals.
func ValidateOpaqueToken(name, value string, opts OpaqueTokenOptions) *Violation {
	if value == "" {
		if opts.AllowEmpty {
			return nil
		}

		return violation(ReasonEmptyValue, name+" must be non-empty")
	}

	if opts.MaxLength > 0 && len(value) > opts.MaxLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("%s length must be <= %d bytes", name, opts.MaxLength),
		)
	}

	if HasControl(value) {
		return violation(ReasonInvalidCharacter, name+" must not contain control bytes")
	}

	if HasWhitespace(value) {
		return violation(ReasonInvalidCharacter, name+" must not contain whitespace")
	}

	if HasUnsafeScalarChar(value) {
		return violation(ReasonInvalidCharacter, name+" must not contain path separators")
	}

	return nil
}
