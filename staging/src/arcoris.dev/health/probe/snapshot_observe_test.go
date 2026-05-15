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
)

func TestSnapshotIsObserved(t *testing.T) {
	t.Parallel()

	observed := time.Unix(9, 0)
	updated := time.Unix(10, 0)

	tests := []struct {
		name string
		snap Snapshot
		want bool
	}{
		{
			name: "zero",
			want: false,
		},
		{
			name: "updated without revision",
			snap: Snapshot{
				Updated: updated,
			},
			want: false,
		},
		{
			name: "revision without updated",
			snap: Snapshot{
				Revision: 1,
			},
			want: false,
		},
		{
			name: "updated and revision without target and report",
			snap: Snapshot{
				Updated:  updated,
				Revision: 1,
			},
			want: false,
		},
		{
			name: "target mismatch is not observed",
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
			want: false,
		},
		{
			name: "invalid report is not observed",
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
			want: false,
		},
		{
			name: "unknown report can still be observed",
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
			name: "stale report can still be observed",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
				Stale:    true,
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.snap.IsObserved(); got != tc.want {
				t.Fatalf("IsObserved() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSnapshotIsFresh(t *testing.T) {
	t.Parallel()

	observed := time.Unix(9, 0)
	updated := time.Unix(10, 0)

	tests := []struct {
		name string
		snap Snapshot
		want bool
	}{
		{
			name: "zero",
			want: false,
		},
		{
			name: "observed fresh",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
			},
			want: true,
		},
		{
			name: "observed stale",
			snap: Snapshot{
				Target: health.TargetReady,
				Report: health.Report{
					Target:   health.TargetReady,
					Status:   health.StatusHealthy,
					Observed: observed,
				},
				Updated:  updated,
				Revision: 1,
				Stale:    true,
			},
			want: false,
		},
		{
			name: "updated and revision without target and report",
			snap: Snapshot{
				Updated:  updated,
				Revision: 1,
			},
			want: false,
		},
		{
			name: "stale without observation",
			snap: Snapshot{
				Stale: true,
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.snap.IsFresh(); got != tc.want {
				t.Fatalf("IsFresh() = %v, want %v", got, tc.want)
			}
		})
	}
}
