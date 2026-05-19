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

// CapacitySnapshot describes the current permit capacity of a Limiter.
//
// Capacity is point-in-time diagnostic state. It is derived from the coherent
// limiter state protected by the limiter mutex and should not be used as a
// substitute for acquiring a Permit.
type CapacitySnapshot struct {
	// Limit is the maximum number of concurrently held permits.
	Limit uint64

	// InFlight is the number of currently held permits.
	InFlight uint64

	// Available is the number of permits that could be acquired at the time this
	// snapshot was built.
	Available uint64

	// Full reports whether no permit capacity was available.
	Full bool
}

// newCapacitySnapshot builds a capacity snapshot from coherent limiter state.
func newCapacitySnapshot(limit, inFlight uint64) CapacitySnapshot {
	available := uint64(0)
	if inFlight < limit {
		available = limit - inFlight
	}

	return CapacitySnapshot{
		Limit:     limit,
		InFlight:  inFlight,
		Available: available,
		Full:      available == 0,
	}
}

// IsValid reports whether s satisfies the bulkhead capacity invariants.
func (s CapacitySnapshot) IsValid() bool {
	if s.Limit == 0 {
		return false
	}
	if s.InFlight > s.Limit {
		return false
	}
	if s.Available != s.Limit-s.InFlight {
		return false
	}

	return s.Full == (s.Available == 0)
}
