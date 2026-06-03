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
	"fmt"
	"strings"
)

// ErrorList contains collected validation diagnostics.
//
// Validation collects multiple independent failures up to Options.MaxErrors.
// Unwrap returns the underlying diagnostics so errors.Is and errors.As can
// inspect every collected error.
type ErrorList []error

// Error summarizes the collected validation diagnostics.
func (l ErrorList) Error() string {
	switch len(l) {
	case 0:
		return "valuevalidation: no errors"
	case 1:
		return l[0].Error()
	default:
		return fmt.Sprintf("valuevalidation: %d errors; first: %v", len(l), l[0])
	}
}

// Unwrap returns the collected errors for errors.Is and errors.As.
func (l ErrorList) Unwrap() []error {
	return l.Errors()
}

// Errors returns the collected diagnostics as a caller-owned slice.
func (l ErrorList) Errors() []error {
	if len(l) == 0 {
		return nil
	}

	out := make([]error, len(l))
	copy(out, l)
	return out
}

// First returns the first diagnostic, or nil when the list is empty.
func (l ErrorList) First() error {
	if len(l) == 0 {
		return nil
	}

	return l[0]
}

// FormatAll renders every collected diagnostic in order.
//
// Error keeps a compact summary for the error interface. FormatAll is for
// callers that intentionally want the full ordered diagnostic list without
// depending on ErrorList's slice representation.
func (l ErrorList) FormatAll() string {
	if len(l) == 0 {
		return ""
	}

	var builder strings.Builder
	for i, err := range l {
		if i > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(err.Error())
	}

	return builder.String()
}

// Len reports the number of collected diagnostics.
func (l ErrorList) Len() int {
	return len(l)
}

// IsEmpty reports whether no diagnostics were collected.
func (l ErrorList) IsEmpty() bool {
	return len(l) == 0
}
