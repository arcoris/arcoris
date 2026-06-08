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

package deadline

import "errors"

var (
	// ErrNilContext reports a nil context at a public deadline API boundary.
	ErrNilContext = errors.New("deadline: nil context")

	// ErrNegativeDuration classifies negative duration inputs at public deadline
	// API boundaries.
	ErrNegativeDuration = errors.New("deadline: negative duration")
)

// NegativeDurationError reports which duration argument was negative.
//
// The error classifies with ErrNegativeDuration. It is used as a panic value
// because negative durations are programmer or configuration errors in this
// package, not ordinary denied start decisions.
type NegativeDurationError struct {
	// Name is the public argument name that received a negative duration.
	Name string
}

// Error returns a stable diagnostic message for e.
func (e NegativeDurationError) Error() string {
	if e.Name == "" {
		return ErrNegativeDuration.Error()
	}
	return "deadline: negative " + e.Name
}

// Is reports whether target is ErrNegativeDuration.
func (e NegativeDurationError) Is(target error) bool {
	return target == ErrNegativeDuration
}
