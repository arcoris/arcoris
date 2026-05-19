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

package bulkhead

// StatsSnapshot describes lifetime admission counters for one Limiter.
//
// Counters are part of the limiter snapshot so diagnostic reads observe capacity
// and accounting from the same coherent state. The base implementation keeps
// these counters under the limiter mutex rather than using independent atomics;
// this preserves exact relationships such as Acquired-Released == InFlight.
type StatsSnapshot struct {
	// Acquired is the lifetime number of successful permit acquisitions.
	Acquired uint64

	// Rejected is the lifetime number of rejected acquisition attempts.
	Rejected uint64

	// Released is the lifetime number of permits released back to the limiter.
	Released uint64
}

// IsValid reports whether s satisfies basic lifetime counter invariants.
func (s StatsSnapshot) IsValid() bool {
	return s.Released <= s.Acquired
}

// InFlight returns the number of currently held permits implied by the lifetime
// counters.
//
// InFlight returns zero for invalid snapshots where Released exceeds Acquired.
func (s StatsSnapshot) InFlight() uint64 {
	if s.Released > s.Acquired {
		return 0
	}

	return s.Acquired - s.Released
}
