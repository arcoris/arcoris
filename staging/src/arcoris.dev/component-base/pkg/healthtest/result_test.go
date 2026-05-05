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

package healthtest

import (
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

func TestResultFixtures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     health.Result
		wantStatus health.Status
		wantReason health.Reason
	}{
		{"healthy", HealthyResult("storage"), health.StatusHealthy, health.ReasonNone},
		{"starting", StartingResult("startup"), health.StatusStarting, health.ReasonStarting},
		{"degraded", DegradedResult("queue", health.ReasonOverloaded), health.StatusDegraded, health.ReasonOverloaded},
		{"unhealthy", UnhealthyResult("database", health.ReasonFatal), health.StatusUnhealthy, health.ReasonFatal},
		{"unknown", UnknownResult("cache", health.ReasonNotObserved), health.StatusUnknown, health.ReasonNotObserved},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			AssertResultStatus(t, tc.result, tc.wantStatus)
			AssertResultReason(t, tc.result, tc.wantReason)
			if !tc.result.Observed.Equal(ObservedTime) {
				t.Fatalf("Observed = %v, want %v", tc.result.Observed, ObservedTime)
			}
		})
	}
}

func TestResultTransformers(t *testing.T) {
	t.Parallel()

	result := ResultWithPrivateCause(ResultWithDuration(HealthyResult("storage"), 15*time.Millisecond))

	if result.Duration != 15*time.Millisecond {
		t.Fatalf("Duration = %s, want 15ms", result.Duration)
	}
	if result.Cause == nil || result.Cause.Error() != "private cause" {
		t.Fatalf("Cause = %v, want private cause", result.Cause)
	}
}
