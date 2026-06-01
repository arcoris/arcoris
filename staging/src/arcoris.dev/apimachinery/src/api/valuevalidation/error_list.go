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

import "fmt"

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
	if len(l) == 0 {
		return nil
	}

	out := make([]error, len(l))
	copy(out, l)
	return out
}

// Len reports the number of collected diagnostics.
func (l ErrorList) Len() int {
	return len(l)
}

// IsEmpty reports whether no diagnostics were collected.
func (l ErrorList) IsEmpty() bool {
	return len(l) == 0
}
