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
	"context"
	"testing"

	"arcoris.dev/health"
)

func TestStaticAndFuncChecker(t *testing.T) {
	t.Parallel()

	static := StaticChecker("storage", HealthyResult("storage"))
	if static.Name() != "storage" {
		t.Fatalf("Name = %q, want storage", static.Name())
	}
	AssertResultStatus(t, static.Check(context.Background()), health.StatusHealthy)

	fn := FuncChecker("queue", func(context.Context) health.Result {
		return DegradedResult("queue", health.ReasonOverloaded)
	})
	AssertResultStatus(t, fn.Check(context.Background()), health.StatusDegraded)
}

func TestStatusCheckers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checker    health.Checker
		wantStatus health.Status
		wantReason health.Reason
	}{
		{"healthy", HealthyChecker("storage"), health.StatusHealthy, health.ReasonNone},
		{"starting", StartingChecker("startup"), health.StatusStarting, health.ReasonStarting},
		{"degraded", DegradedChecker("queue", health.ReasonOverloaded), health.StatusDegraded, health.ReasonOverloaded},
		{"unhealthy", UnhealthyChecker("database", health.ReasonFatal), health.StatusUnhealthy, health.ReasonFatal},
		{"unknown", UnknownChecker("cache", health.ReasonNotObserved), health.StatusUnknown, health.ReasonNotObserved},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := tc.checker.Check(context.Background())
			AssertResultStatus(t, result, tc.wantStatus)
			AssertResultReason(t, result, tc.wantReason)
			if result.Name != tc.checker.Name() {
				t.Fatalf("result name = %q, want checker name %q", result.Name, tc.checker.Name())
			}
		})
	}
}

func TestCheckerAllowsInvalidNames(t *testing.T) {
	t.Parallel()

	checker := HealthyChecker("invalid name")
	if checker.Name() != "invalid name" {
		t.Fatalf("Name = %q, want invalid name", checker.Name())
	}
}
