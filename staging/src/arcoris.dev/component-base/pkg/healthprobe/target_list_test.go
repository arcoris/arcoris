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
	"errors"
	"testing"

	"arcoris.dev/component-base/pkg/health"
)

func TestNormalizeTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		targets []health.Target
		want    []health.Target
	}{
		{
			name:    "single target",
			targets: []health.Target{health.TargetReady},
			want:    []health.Target{health.TargetReady},
		},
		{
			name: "all concrete targets in caller order",
			targets: []health.Target{
				health.TargetReady,
				health.TargetLive,
				health.TargetStartup,
			},
			want: []health.Target{
				health.TargetReady,
				health.TargetLive,
				health.TargetStartup,
			},
		},
		{
			name: "all built in concrete targets",
			targets: []health.Target{
				health.TargetStartup,
				health.TargetLive,
				health.TargetReady,
			},
			want: []health.Target{
				health.TargetStartup,
				health.TargetLive,
				health.TargetReady,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeTargets(tc.targets)
			if err != nil {
				t.Fatalf("normalizeTargets(%v) = %v, want nil", tc.targets, err)
			}

			if !sameTargets(got, tc.want) {
				t.Fatalf("normalizeTargets(%v) = %v, want %v", tc.targets, got, tc.want)
			}
		})
	}
}

func TestNormalizeTargetsRejectsEmptyList(t *testing.T) {
	t.Parallel()

	_, err := normalizeTargets(nil)

	if !errors.Is(err, ErrNoTargets) {
		t.Fatalf("normalizeTargets(nil) = %v, want ErrNoTargets", err)
	}
}

func TestNormalizeTargetsRejectsNonConcreteTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target health.Target
	}{
		{
			name:   "unknown",
			target: health.TargetUnknown,
		},
		{
			name:   "invalid",
			target: health.Target(255),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := normalizeTargets([]health.Target{tc.target})

			if !errors.Is(err, health.ErrInvalidTarget) {
				t.Fatalf(
					"normalizeTargets(%v) = %v, want health.ErrInvalidTarget",
					tc.target,
					err,
				)
			}

			var targetErr health.InvalidTargetError
			if !errors.As(err, &targetErr) {
				t.Fatalf("errors.As(%T, health.InvalidTargetError) = false, want true", err)
			}
			if targetErr.Target != tc.target {
				t.Fatalf("InvalidTargetError.Target = %s, want %s", targetErr.Target, tc.target)
			}
		})
	}
}

func TestNormalizeTargetsRejectsDuplicateTarget(t *testing.T) {
	t.Parallel()

	targets := []health.Target{
		health.TargetReady,
		health.TargetLive,
		health.TargetReady,
	}

	_, err := normalizeTargets(targets)

	if !errors.Is(err, ErrDuplicateTarget) {
		t.Fatalf("normalizeTargets(%v) = %v, want ErrDuplicateTarget", targets, err)
	}

	var duplicateErr DuplicateTargetError
	if !errors.As(err, &duplicateErr) {
		t.Fatalf("errors.As(%T, DuplicateTargetError) = false, want true", err)
	}
	if duplicateErr.Target != health.TargetReady {
		t.Fatalf("Target = %s, want %s", duplicateErr.Target, health.TargetReady)
	}
	if duplicateErr.Index != 2 {
		t.Fatalf("Index = %d, want 2", duplicateErr.Index)
	}
	if duplicateErr.PreviousIndex != 0 {
		t.Fatalf("PreviousIndex = %d, want 0", duplicateErr.PreviousIndex)
	}
}

func TestNormalizeTargetsReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	source := []health.Target{
		health.TargetReady,
		health.TargetLive,
	}

	got, err := normalizeTargets(source)
	if err != nil {
		t.Fatalf("normalizeTargets(%v) = %v, want nil", source, err)
	}

	got[0] = health.TargetStartup

	if source[0] != health.TargetReady {
		t.Fatalf("source was mutated through result: source[0]=%s", source[0])
	}
}

func TestCopyTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		targets []health.Target
		want    []health.Target
	}{
		{
			name:    "nil",
			targets: nil,
			want:    nil,
		},
		{
			name:    "empty",
			targets: []health.Target{},
			want:    nil,
		},
		{
			name: "values",
			targets: []health.Target{
				health.TargetLive,
				health.TargetReady,
			},
			want: []health.Target{
				health.TargetLive,
				health.TargetReady,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := copyTargets(tc.targets)
			if !sameTargets(got, tc.want) {
				t.Fatalf("copyTargets(%v) = %v, want %v", tc.targets, got, tc.want)
			}

			if len(got) > 0 {
				got[0] = health.TargetStartup
				if tc.targets[0] == health.TargetStartup {
					t.Fatal("copyTargets returned slice sharing caller backing array")
				}
			}
		})
	}
}

func TestContainsTarget(t *testing.T) {
	t.Parallel()

	targets := []health.Target{
		health.TargetStartup,
		health.TargetReady,
	}

	tests := []struct {
		name   string
		target health.Target
		want   bool
	}{
		{
			name:   "contains startup",
			target: health.TargetStartup,
			want:   true,
		},
		{
			name:   "contains ready",
			target: health.TargetReady,
			want:   true,
		},
		{
			name:   "does not contain live",
			target: health.TargetLive,
			want:   false,
		},
		{
			name:   "does not contain unknown",
			target: health.TargetUnknown,
			want:   false,
		},
		{
			name:   "does not contain invalid",
			target: health.Target(255),
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := containsTarget(targets, tc.target); got != tc.want {
				t.Fatalf("containsTarget(%v, %s) = %v, want %v", targets, tc.target, got, tc.want)
			}
		})
	}
}
