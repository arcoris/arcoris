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

// Snapshot is the domain read model published by a Limiter.
//
// Snapshot is intentionally separate from snapshot.Snapshot[Snapshot]. This type
// contains bulkhead state only; source-local revision and update timestamps are
// provided by package snapshot.
type Snapshot struct {
	// Capacity describes current permit usage and available capacity.
	Capacity CapacitySnapshot

	// Stats describes lifetime admission counters.
	Stats StatsSnapshot
}

// IsValid reports whether s satisfies the bulkhead snapshot invariants.
func (s Snapshot) IsValid() bool {
	if !s.Capacity.IsValid() || !s.Stats.IsValid() {
		return false
	}

	return s.Stats.InFlight() == s.Capacity.InFlight
}
