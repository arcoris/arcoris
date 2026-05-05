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

package healthprobe

import (
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

func TestSnapshotIsValid(t *testing.T) {
	t.Parallel()

	observed := time.Unix(9, 0)
	updated := time.Unix(10, 0)

	tests := []struct {
		name     string
		snapshot Snapshot
		want     bool
	}{
		{
			name: "zero snapshot is valid absence of observation",
			want: true,
		},
		{
			name: "valid healthy snapshot",
			snapshot: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
					Checks: []health.Result{
						health.Healthy("database").WithObserved(observed),
					},
				},
				Updated:    updated,
				Generation: 1,
			},
			want: true,
		},
		{
			name: "valid unknown snapshot",
			snapshot: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusUnknown,
					Observed: observed,
				},
				Updated:    updated,
				Generation: 1,
			},
			want: true,
		},
		{
			name: "valid stale snapshot",
			snapshot: Snapshot{
				Target: health.TargetLive,
				Report: health.Report{
					Target:   health.TargetLive,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:    updated,
				Generation: 2,
				Stale:      true,
			},
			want: true,
		},
		{
			name: "updated and generation without target and report",
			snapshot: Snapshot{
				Updated:    updated,
				Generation: 1,
			},
		},
		{
			name: "non concrete target",
			snapshot: Snapshot{
				Target: health.TargetUnknown,
				Report: health.Report{
					Target:   health.TargetUnknown,
					Status:   health.StatusUnknown,
					Observed: observed,
				},
				Updated:    updated,
				Generation: 1,
			},
		},
		{
			name: "target mismatch",
			snapshot: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetLive,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:    updated,
				Generation: 1,
			},
		},
		{
			name: "invalid report status",
			snapshot: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.Status(255),
					Observed: observed,
				},
				Updated:    updated,
				Generation: 1,
			},
		},
		{
			name: "invalid report duration",
			snapshot: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
					Duration: -time.Millisecond,
				},
				Updated:    updated,
				Generation: 1,
			},
		},
		{
			name: "missing updated",
			snapshot: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Generation: 1,
			},
		},
		{
			name: "missing generation",
			snapshot: Snapshot{
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
			snapshot: Snapshot{
				Stale: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.snapshot.IsValid(); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}
