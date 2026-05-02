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
		result Result
		status Status
		reason Reason
	}{
		{"healthy", Healthy("storage"), StatusHealthy, ReasonNone},
		{"starting", Starting("storage", ReasonStarting, "starting"), StatusStarting, ReasonStarting},
		{"degraded", Degraded("storage", ReasonOverloaded, "overloaded"), StatusDegraded, ReasonOverloaded},
		{"unhealthy", Unhealthy("storage", ReasonFatal, "fatal"), StatusUnhealthy, ReasonFatal},
		{"unknown", Unknown("storage", ReasonNotObserved, "unknown"), StatusUnknown, ReasonNotObserved},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.result.Name != "storage" {
				t.Fatalf("Name = %q, want storage", test.result.Name)
			}
			if test.result.Status != test.status {
				t.Fatalf("Status = %s, want %s", test.result.Status, test.status)
			}
			if test.result.Reason != test.reason {
				t.Fatalf("Reason = %s, want %s", test.result.Reason, test.reason)
			}
		})
	}
}

func TestResultWithMethods(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")
	observed := testObserved
	result := Healthy("storage").
		WithCause(cause).
		WithObserved(observed).
		WithDuration(time.Second).
		WithMessage("message").
		WithReason(ReasonFatal)

	if result.Cause != cause || result.Observed != observed || result.Duration != time.Second {
		t.Fatalf("metadata not preserved: %+v", result)
	}
	if result.Message != "message" || result.Reason != ReasonFatal {
		t.Fatalf("message/reason not preserved: %+v", result)
	}
}

func TestResultNormalize(t *testing.T) {
	t.Parallel()

	result := Result{Status: Status(99), Duration: -time.Second}.Normalize("storage", testObserved)

	if result.Name != "storage" {
		t.Fatalf("Name = %q, want storage", result.Name)
	}
	if result.Status != StatusUnknown {
		t.Fatalf("Status = %s, want unknown", result.Status)
	}
	if result.Observed != testObserved {
		t.Fatalf("Observed = %v, want %v", result.Observed, testObserved)
	}
	if result.Duration != 0 {
		t.Fatalf("Duration = %s, want 0", result.Duration)
	}
}

func TestResultPredicates(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")
	result := Degraded("storage", ReasonOverloaded, "overloaded").
		WithCause(cause).
		WithObserved(testObserved)

	if !result.IsValid() || !result.IsNamed() || !result.IsObserved() || !result.HasCause() {
		t.Fatalf("predicate mismatch for %+v", result)
	}
	if !result.HasReason(ReasonOverloaded) || result.HasReason(ReasonFatal) {
		t.Fatalf("reason predicate mismatch for %+v", result)
	}
	if result.IsAffirmative() || result.IsNegative() || !result.IsKnown() || !result.IsOperational() {
		t.Fatalf("status predicate mismatch for %+v", result)
	}
	if !Unhealthy("storage", ReasonFatal, "fatal").MoreSevereThan(result) {
		t.Fatal("unhealthy result should be more severe than degraded")
	}
	if (Result{Status: StatusHealthy, Duration: -time.Second}).IsValid() {
		t.Fatal("negative duration result should be invalid")
	}
}
