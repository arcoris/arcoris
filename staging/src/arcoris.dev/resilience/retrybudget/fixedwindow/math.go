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

package fixedwindow

import "math"

// allowedRetries returns the retry capacity for original attempts.
//
// The formula is:
//
//	min + floor(original * ratio)
//
// The result saturates at math.MaxUint64 instead of wrapping.
func allowedRetries(original uint64, ratio float64, min uint64) uint64 {
	if ratio <= 0 || original == 0 {
		return min
	}

	extra := math.Floor(float64(original) * ratio)
	if extra <= 0 {
		return min
	}

	remaining := math.MaxUint64 - min
	if extra >= float64(remaining) {
		return math.MaxUint64
	}

	return min + uint64(extra)
}

// availableRetries returns allowed - used, saturating at zero.
func availableRetries(allowed, used uint64) uint64 {
	if used >= allowed {
		return 0
	}
	return allowed - used
}

// saturatingInc returns n+1 or math.MaxUint64 when n is already saturated.
func saturatingInc(n uint64) uint64 {
	if n == math.MaxUint64 {
		return math.MaxUint64
	}
	return n + 1
}
