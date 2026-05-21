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

package probe

import (
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/snapshot"
)

func TestSnapshotObservedRequiresRevision(t *testing.T) {
	t.Parallel()

	observed := time.Unix(9, 0)
	snap := Snapshot{
		Target: health.TargetReady,
		Report: health.Report{
			Target:   health.TargetReady,
			Status:   health.StatusHealthy,
			Observed: observed,
		},
		Revision: snapshot.ZeroRevision,
		Updated:  time.Unix(10, 0),
	}

	if snap.IsValid() {
		t.Fatal("IsValid() = true, want false")
	}
	if snap.IsObserved() {
		t.Fatal("IsObserved() = true, want false")
	}
}

func TestSnapshotIsValid(t *testing.T) {
	t.Parallel()

	observed := time.Unix(9, 0)
	updated := time.Unix(10, 0)

	tests := []struct {
		name string
		snap Snapshot
		want bool
	}{
		{
			name: "zero snapshot is valid absence of observation",
			want: true,
		},
		{
			name: "valid healthy snapshot",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
					Checks: []health.Result{
						health.Healthy("database").WithObserved(observed),
					},
				},
				Updated:  updated,
				Revision: 1,
			},
			want: true,
		},
		{
			name: "valid unknown snapshot",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusUnknown,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
			},
			want: true,
		},
		{
			name: "valid stale snapshot",
			snap: Snapshot{
				Target: health.TargetLive,
				Report: health.Report{
					Target:   health.TargetLive,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 2,
				Stale:    true,
			},
			want: true,
		},
		{
			name: "updated and revision without target and report",
			snap: Snapshot{
				Updated:  updated,
				Revision: 1,
			},
		},
		{
			name: "non concrete target",
			snap: Snapshot{
				Target: health.TargetUnknown,
				Report: health.Report{
					Target:   health.TargetUnknown,
					Status:   health.StatusUnknown,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
			},
		},
		{
			name: "target mismatch",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetLive,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
			},
		},
		{
			name: "invalid report status",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.Status(255),
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
			},
		},
		{
			name: "invalid report duration",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
					Duration: -time.Millisecond,
				},
				Updated:  updated,
				Revision: 1,
			},
		},
		{
			name: "missing updated",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Revision: 1,
			},
		},
		{
			name: "missing revision",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated: updated,
			},
		},
		{
			name: "stale without observation",
			snap: Snapshot{
				Stale: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.snap.IsValid(); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}
