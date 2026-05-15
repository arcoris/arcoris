/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package retrybudget

// CapacitySnapshot reports the retry admission capacity currently exposed by an
// implementation.
type CapacitySnapshot struct {
	// Allowed is the total retry capacity currently allowed by the implementation.
	Allowed uint64

	// Available is the remaining retry capacity currently available for admission.
	Available uint64

	// Exhausted reports whether no retry capacity remains.
	Exhausted bool
}

// IsValid reports whether s is internally consistent.
func (s CapacitySnapshot) IsValid() bool {
	if s.Available > s.Allowed {
		return false
	}
	if s.Exhausted {
		return s.Available == 0
	}
	return s.Available > 0
}

// HasAvailable reports whether at least one retry attempt may still be admitted
// according to the published capacity.
func (s CapacitySnapshot) HasAvailable() bool {
	return !s.Exhausted && s.Available > 0
}
