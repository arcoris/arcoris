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

func TestSnapshotZeroUsesZeroRevision(t *testing.T) {
	t.Parallel()

	var snap Snapshot

	if !snap.IsValid() {
		t.Fatal("zero Snapshot IsValid() = false, want true")
	}
	if snap.Revision != snapshot.ZeroRevision {
		t.Fatalf("zero Snapshot Revision = %d, want zero", snap.Revision)
	}
}

func TestSnapshotIsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		snapshot Snapshot
		want     bool
	}{
		{
			name: "zero",
			want: true,
		},
		{
			name: "target set",
			snapshot: Snapshot{
				Target: health.TargetReady,
			},
		},
		{
			name: "report target set",
			snapshot: Snapshot{
				Report: health.Report{
					Target: health.TargetReady,
					Status: health.StatusUnknown,
				},
			},
		},
		{
			name: "report status set",
			snapshot: Snapshot{
				Report: health.Report{
					Target: health.TargetUnknown,
					Status: health.StatusHealthy,
				},
			},
		},
		{
			name: "report observed set",
			snapshot: Snapshot{
				Report: health.Report{
					Target:   health.TargetUnknown,
					Status:   health.StatusUnknown,
					Observed: time.Unix(10, 0),
				},
			},
		},
		{
			name: "report duration set",
			snapshot: Snapshot{
				Report: health.Report{
					Target:   health.TargetUnknown,
					Status:   health.StatusUnknown,
					Duration: time.Millisecond,
				},
			},
		},
		{
			name: "report checks set",
			snapshot: Snapshot{
				Report: health.Report{
					Target: health.TargetUnknown,
					Status: health.StatusUnknown,
					Checks: []health.Result{
						health.Healthy("database"),
					},
				},
			},
		},
		{
			name: "updated set",
			snapshot: Snapshot{
				Updated: time.Unix(10, 0),
			},
		},
		{
			name: "revision set",
			snapshot: Snapshot{
				Revision: 1,
			},
		},
		{
			name: "stale set",
			snapshot: Snapshot{
				Stale: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.snapshot.IsZero(); got != tc.want {
				t.Fatalf("IsZero() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestReportIsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		report health.Report
		want   bool
	}{
		{
			name: "zero",
			want: true,
		},
		{
			name: "target set",
			report: health.Report{
				Target: health.TargetReady,
				Status: health.StatusUnknown,
			},
			want: false,
		},
		{
			name: "status set",
			report: health.Report{
				Target: health.TargetUnknown,
				Status: health.StatusHealthy,
			},
			want: false,
		},
		{
			name: "observed set",
			report: health.Report{
				Target:   health.TargetUnknown,
				Status:   health.StatusUnknown,
				Observed: time.Unix(10, 0),
			},
			want: false,
		},
		{
			name: "duration set",
			report: health.Report{
				Target:   health.TargetUnknown,
				Status:   health.StatusUnknown,
				Duration: time.Millisecond,
			},
			want: false,
		},
		{
			name: "checks set",
			report: health.Report{
				Target: health.TargetUnknown,
				Status: health.StatusUnknown,
				Checks: []health.Result{
					health.Healthy("database"),
				},
			},
			want: false,
		},
		{
			name: "nil checks remains zero",
			report: health.Report{
				Target: health.TargetUnknown,
				Status: health.StatusUnknown,
				Checks: nil,
			},
			want: true,
		},
		{
			name: "empty checks remains zero",
			report: health.Report{
				Target: health.TargetUnknown,
				Status: health.StatusUnknown,
				Checks: []health.Result{},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := reportIsZero(tc.report); got != tc.want {
				t.Fatalf("reportIsZero() = %v, want %v", got, tc.want)
			}
		})
	}
}
