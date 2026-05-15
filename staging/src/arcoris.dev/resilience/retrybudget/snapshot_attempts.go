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

import "math"

// AttemptsSnapshot reports attempt counters observed by a retry-budget
// implementation.
type AttemptsSnapshot struct {
	// Original is the number of original, non-retry attempts observed by the
	// implementation in the scope represented by the snapshot.
	Original uint64

	// Retry is the number of retry attempts admitted or observed by the
	// implementation in the scope represented by the snapshot.
	Retry uint64
}

// Total returns Original + Retry, saturating at math.MaxUint64 on overflow.
func (s AttemptsSnapshot) Total() uint64 {
	if math.MaxUint64-s.Original < s.Retry {
		return math.MaxUint64
	}
	return s.Original + s.Retry
}

// HasTraffic reports whether any original or retry attempts were observed.
func (s AttemptsSnapshot) HasTraffic() bool {
	return s.Original != 0 || s.Retry != 0
}
