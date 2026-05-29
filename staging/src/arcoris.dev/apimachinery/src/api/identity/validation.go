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

package identity

import (
	"strings"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// validateDNS1123SingleLabel checks resource-like one-label identities.
func validateDNS1123SingleLabel(name, value string, allowEmpty bool) error {
	if value == "" {
		if allowEmpty {
			return nil
		}
		return invalid(name, value, ErrorReasonEmptyValue, name+" must be non-empty")
	}
	if strings.Contains(value, dnsLabelSeparator) || strings.Contains(value, subresourceSeparator) {
		return invalid(name, value, ErrorReasonInvalidCharacter, name+" must be a DNS-1123 single label")
	}
	if violation := lexical.ValidateDNS1123Label(value); violation != nil {
		return invalid(name, value, errorReasonFromLexical(violation.Reason), violation.Detail)
	}
	return nil
}

// errorReasonFromLexical maps internal lexical reasons to identity reasons.
func errorReasonFromLexical(reason lexical.Reason) ErrorReason {
	switch reason {
	case lexical.ReasonEmptyValue:
		return ErrorReasonEmptyValue
	case lexical.ReasonInvalidLength:
		return ErrorReasonInvalidLength
	case lexical.ReasonInvalidCharacter:
		return ErrorReasonInvalidCharacter
	case lexical.ReasonInvalidEdge:
		return ErrorReasonInvalidEdge
	default:
		return ErrorReasonInvalidForm
	}
}

// isLower reports whether b is an ASCII lowercase letter.
func isLower(b byte) bool { return lexical.IsASCIILower(b) }

// isUpper reports whether b is an ASCII uppercase letter.
func isUpper(b byte) bool { return lexical.IsASCIIUpper(b) }

// isDigit reports whether b is an ASCII digit.
func isDigit(b byte) bool { return lexical.IsASCIIDigit(b) }
