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

// Reason identifies the metadata grammar invariant that failed.
type Reason string

// Reasons classify internal metadata grammar failures for domain packages.
const (
	// ReasonEmptyValue reports a required metadata token that is absent.
	ReasonEmptyValue Reason = "empty_value"
	// ReasonInvalidLength reports a byte-length limit violation.
	ReasonInvalidLength Reason = "invalid_length"
	// ReasonInvalidCharacter reports a byte outside the grammar.
	ReasonInvalidCharacter Reason = "invalid_character"
	// ReasonInvalidEdge reports an invalid first or last byte.
	ReasonInvalidEdge Reason = "invalid_edge"
	// ReasonInvalidForm reports a malformed composite token.
	ReasonInvalidForm Reason = "invalid_form"
)

// Violation describes one internal metadata grammar failure.
type Violation struct {
	// Reason is the stable machine-readable grammar failure class.
	Reason Reason
	// Detail explains the concrete byte, length, or form problem.
	Detail string
}

// Error returns a human-readable violation diagnostic.
func (v *Violation) Error() string {
	if v == nil {
		return "<nil>"
	}
	if v.Detail == "" {
		return string(v.Reason)
	}
	return fmt.Sprintf("%s: %s", v.Reason, v.Detail)
}

// violation constructs one metadata grammar failure.
func violation(reason Reason, detail string) *Violation {
	return &Violation{Reason: reason, Detail: detail}
}

// fromLexical maps shared lexical failures to metadata grammar failures.
func fromLexical(v *lexical.Violation) *Violation {
	if v == nil {
		return nil
	}
	return violation(reasonFromLexical(v.Reason), v.Detail)
}

// reasonFromLexical preserves lexical diagnostic precision without exporting
// lexical.Reason through public metadata APIs.
func reasonFromLexical(reason lexical.Reason) Reason {
	switch reason {
	case lexical.ReasonEmptyValue:
		return ReasonEmptyValue
	case lexical.ReasonInvalidLength:
		return ReasonInvalidLength
	case lexical.ReasonInvalidCharacter:
		return ReasonInvalidCharacter
	case lexical.ReasonInvalidEdge:
		return ReasonInvalidEdge
	default:
		return ReasonInvalidForm
	}
}
