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

package lexical

import "strings"

// Reason identifies the low-level lexical invariant that failed.
//
// Reasons are intentionally small and domain-neutral. Public API packages map
// them into their own ErrorReason values before returning diagnostics.
type Reason string

const (
	// ReasonEmptyValue reports a required token that was absent.
	ReasonEmptyValue Reason = "empty_value"

	// ReasonInvalidLength reports a byte-length limit violation.
	ReasonInvalidLength Reason = "invalid_length"

	// ReasonInvalidCharacter reports a byte outside the accepted grammar.
	ReasonInvalidCharacter Reason = "invalid_character"

	// ReasonInvalidEdge reports an invalid first or last byte.
	ReasonInvalidEdge Reason = "invalid_edge"

	// ReasonInvalidForm reports a structurally invalid token arrangement.
	ReasonInvalidForm Reason = "invalid_form"
)

// Violation describes one failed lexical check.
//
// A nil *Violation means validation succeeded. The type is internal
// implementation data; public API packages should translate it into their own
// domain-specific errors.
type Violation struct {
	// Reason is the stable low-level failure category.
	Reason Reason

	// Detail explains the failure in human-facing terms.
	Detail string
}

// Error returns a compact diagnostic for tests and internal debugging.
func (v Violation) Error() string {
	parts := make([]string, 0, 2)
	if v.Reason != "" {
		parts = append(parts, string(v.Reason))
	}
	if v.Detail != "" {
		parts = append(parts, v.Detail)
	}
	return strings.Join(parts, ": ")
}

// violation creates a structured lexical failure.
func violation(reason Reason, detail string) *Violation {
	return &Violation{Reason: reason, Detail: detail}
}
