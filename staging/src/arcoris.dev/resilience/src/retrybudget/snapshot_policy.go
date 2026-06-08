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

package retrybudget

// PolicySnapshot reports public retry-budget policy parameters when an
// implementation exposes ratio/minimum style capacity settings.
type PolicySnapshot struct {
	// Ratio is the exact retry allowance ratio exposed by the implementation.
	Ratio Ratio

	// Minimum is the minimum retry allowance exposed by the implementation.
	Minimum uint64

	// Bounded reports whether Ratio and Minimum are meaningful.
	Bounded bool
}

// IsValid reports whether s is internally consistent.
func (s PolicySnapshot) IsValid() bool {
	if !s.Bounded {
		return s.Ratio == (Ratio{}) && s.Minimum == 0
	}
	return s.Ratio.IsValid()
}

// IsBounded reports whether s exposes finite retry-budget policy parameters.
func (s PolicySnapshot) IsBounded() bool {
	return s.Bounded
}

// HasMinimum reports whether the policy exposes a positive minimum retry
// allowance.
func (s PolicySnapshot) HasMinimum() bool {
	return s.Bounded && s.Minimum > 0
}
