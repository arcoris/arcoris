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

import "testing"

func TestStatsSnapshotIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		snap StatsSnapshot
		want bool
	}{
		{name: "zero", snap: StatsSnapshot{}, want: true},
		{name: "balanced", snap: StatsSnapshot{Acquired: 3, Released: 3}, want: true},
		{name: "in flight", snap: StatsSnapshot{Acquired: 3, Released: 1}, want: true},
		{name: "released exceeds acquired", snap: StatsSnapshot{Acquired: 1, Released: 2}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.snap.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsSnapshotInFlight(t *testing.T) {
	t.Parallel()

	if got := (StatsSnapshot{Acquired: 5, Released: 2}).InFlight(); got != 3 {
		t.Fatalf("InFlight() = %d, want 3", got)
	}
	if got := (StatsSnapshot{Acquired: 1, Released: 2}).InFlight(); got != 0 {
		t.Fatalf("invalid InFlight() = %d, want 0", got)
	}
}
