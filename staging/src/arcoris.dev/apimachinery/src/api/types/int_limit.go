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

// intLimit stores optional int constraints without pointers.
//
// The set bit distinguishes an explicit zero from an unset limit and avoids
// hidden heap allocations in builder-heavy descriptor construction.
type intLimit struct {
	// value is meaningful only when set is true.
	//
	// Validation code must check set before treating value as an active rule.
	value int
	// set reports whether value was explicitly configured.
	//
	// value=0 with set=false means no configured rule; value=0 with set=true
	// means the rule is explicitly zero.
	set bool
}

// validateIntLimits checks non-negative inclusive structural length limits.
//
// intLimit is reserved for descriptor sizes such as string length, bytes
// length, list length, map length, decimal precision, and decimal scale. It is
// not used for fixed-width numeric value bounds.
func validateIntLimits(min, max intLimit, path string) error {
	if min.set && min.value < 0 {
		return typeError(path, ErrInvalidType)
	}
	if max.set && max.value < 0 {
		return typeError(path, ErrInvalidType)
	}
	if min.set && max.set && min.value > max.value {
		return typeError(path, ErrInvalidType)
	}
	return nil
}
