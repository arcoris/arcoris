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

func TestSnapshotIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		snap Snapshot
		want bool
	}{
		{
			name: "valid",
			snap: Snapshot{
				Capacity: CapacitySnapshot{Limit: 2, InFlight: 1, Available: 1, Full: false},
				Stats:    StatsSnapshot{Acquired: 2, Released: 1},
			},
			want: true,
		},
		{
			name: "capacity invalid",
			snap: Snapshot{
				Capacity: CapacitySnapshot{},
				Stats:    StatsSnapshot{},
			},
			want: false,
		},
		{
			name: "stats invalid",
			snap: Snapshot{
				Capacity: CapacitySnapshot{Limit: 2, InFlight: 0, Available: 2, Full: false},
				Stats:    StatsSnapshot{Acquired: 0, Released: 1},
			},
			want: false,
		},
		{
			name: "in flight mismatch",
			snap: Snapshot{
				Capacity: CapacitySnapshot{Limit: 2, InFlight: 1, Available: 1, Full: false},
				Stats:    StatsSnapshot{Acquired: 2, Released: 2},
			},
			want: false,
		},
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
