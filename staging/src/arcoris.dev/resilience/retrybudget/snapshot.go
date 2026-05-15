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

// Snapshot is the domain state published by a retry budget implementation.
//
// Snapshot intentionally does not contain a revision field. Revisions are owned
// by arcoris.dev/snapshot and are carried by snapshot.Snapshot[Snapshot]. This
// keeps retry-budget state as a plain domain value and preserves one shared
// snapshot mechanism across resilience components.
type Snapshot struct {
	// Kind identifies the implementation that produced this snapshot.
	Kind Kind

	// Attempts reports original and retry attempt counters as observed by the
	// implementation.
	Attempts AttemptsSnapshot

	// Capacity reports current retry admission capacity.
	Capacity CapacitySnapshot

	// Window reports the time window associated with the snapshot, when the
	// implementation is window-based.
	Window WindowSnapshot

	// Policy reports public policy parameters associated with the implementation,
	// when the implementation exposes such parameters.
	Policy PolicySnapshot
}

// IsValid reports whether s is an internally consistent retry-budget snapshot.
func (s Snapshot) IsValid() bool {
	return s.Kind.IsValid() &&
		s.Capacity.IsValid() &&
		s.Window.IsValid() &&
		s.Policy.IsValid()
}

// Exhausted reports whether the snapshot has no retry capacity available.
func (s Snapshot) Exhausted() bool {
	return s.Capacity.Exhausted
}

// HasTraffic reports whether the snapshot contains any observed attempt traffic.
func (s Snapshot) HasTraffic() bool {
	return s.Attempts.HasTraffic()
}
