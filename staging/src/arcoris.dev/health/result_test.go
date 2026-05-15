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

package health

import (
	"errors"
	"testing"
	"time"
)

func TestResultConstructors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		res    Result
		status Status
		reason Reason
	}{
		{"healthy", Healthy("storage"), StatusHealthy, ReasonNone},
		{"starting", Starting("storage", ReasonStarting, "starting"), StatusStarting, ReasonStarting},
		{"degraded", Degraded("storage", ReasonOverloaded, "overloaded"), StatusDegraded, ReasonOverloaded},
		{"unhealthy", Unhealthy("storage", ReasonFatal, "fatal"), StatusUnhealthy, ReasonFatal},
		{"unknown", Unknown("storage", ReasonNotObserved, "unknown"), StatusUnknown, ReasonNotObserved},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.res.Name != "storage" {
				t.Fatalf("Name = %q, want storage", tc.res.Name)
			}
			if tc.res.Status != tc.status {
				t.Fatalf("Status = %s, want %s", tc.res.Status, tc.status)
			}
			if tc.res.Reason != tc.reason {
				t.Fatalf("Reason = %s, want %s", tc.res.Reason, tc.reason)
			}
		})
	}
}

func TestResultWithMethods(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")
	observed := testObserved
	res := Healthy("storage").
		WithCause(cause).
		WithObserved(observed).
		WithDuration(time.Second).
		WithMessage("message").
		WithReason(ReasonFatal)

	if res.Cause != cause || res.Observed != observed || res.Duration != time.Second {
		t.Fatalf("metadata not preserved: %+v", res)
	}
	if res.Message != "message" || res.Reason != ReasonFatal {
		t.Fatalf("message/reason not preserved: %+v", res)
	}
}

func TestResultNormalize(t *testing.T) {
	t.Parallel()

	res := Result{Status: Status(99), Duration: -time.Second}.Normalize("storage", testObserved)

	if res.Name != "storage" {
		t.Fatalf("Name = %q, want storage", res.Name)
	}
	if res.Status != StatusUnknown {
		t.Fatalf("Status = %s, want unknown", res.Status)
	}
	if res.Observed != testObserved {
		t.Fatalf("Observed = %v, want %v", res.Observed, testObserved)
	}
	if res.Duration != 0 {
		t.Fatalf("Duration = %s, want 0", res.Duration)
	}
}

func TestResultPredicates(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")
	res := Degraded("storage", ReasonOverloaded, "overloaded").
		WithCause(cause).
		WithObserved(testObserved)

	if !res.IsValid() || !res.IsNamed() || !res.IsObserved() || !res.HasCause() {
		t.Fatalf("predicate mismatch for %+v", res)
	}
	if !res.HasReason(ReasonOverloaded) || res.HasReason(ReasonFatal) {
		t.Fatalf("reason predicate mismatch for %+v", res)
	}
	if res.IsAffirmative() || res.IsNegative() || !res.IsKnown() || !res.IsOperational() {
		t.Fatalf("status predicate mismatch for %+v", res)
	}
	if !Unhealthy("storage", ReasonFatal, "fatal").MoreSevereThan(res) {
		t.Fatal("unhealthy result should be more severe than degraded")
	}
	if (Result{Status: StatusHealthy, Duration: -time.Second}).IsValid() {
		t.Fatal("negative duration result should be invalid")
	}
	if (Result{Status: StatusHealthy, Reason: Reason("bad-reason")}).IsValid() {
		t.Fatal("invalid reason result should be invalid")
	}
}
