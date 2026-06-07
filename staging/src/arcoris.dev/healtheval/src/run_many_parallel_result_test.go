// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eval

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/healthregistry"
)

func TestEvaluatorParallelAggregatesMostSevereStatus(t *testing.T) {
	t.Parallel()

	registry := healthregistry.NewBuilder()
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "healthy", func(context.Context) health.Result {
		return health.Healthy("healthy")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "degraded", func(context.Context) health.Result {
		return health.Degraded("degraded", health.ReasonOverloaded, "degraded")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "unhealthy", func(context.Context) health.Result {
		return health.Unhealthy("unhealthy", health.ReasonFatal, "unhealthy")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, 3),
	)

	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Status != health.StatusUnhealthy {
		t.Fatalf("health.Status = %s, want unhealthy", report.Status)
	}
}

func TestEvaluatorParallelPreservesNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkName  string
		fn         health.CheckFunc
		wantStatus health.Status
		wantReason health.Reason
	}{
		{
			name:       "panic",
			checkName:  "panic_check",
			fn:         func(context.Context) health.Result { panic("boom") },
			wantStatus: health.StatusUnhealthy,
			wantReason: health.ReasonPanic,
		},
		{
			name:      "invalid reason",
			checkName: "invalid_reason",
			fn: func(context.Context) health.Result {
				return health.Unknown("invalid_reason", health.Reason("bad-reason"), "bad")
			},
			wantStatus: health.StatusUnknown,
			wantReason: health.ReasonMisconfigured,
		},
		{
			name:       "mismatched name",
			checkName:  "mismatched_name",
			fn:         func(context.Context) health.Result { return health.Healthy("other_name") },
			wantStatus: health.StatusUnknown,
			wantReason: health.ReasonMisconfigured,
		},
		{
			name:      "cause preserved internally",
			checkName: "cause_check",
			fn: func(context.Context) health.Result {
				return health.Unhealthy("cause_check", health.ReasonFatal, "failed").WithCause(errors.New("private"))
			},
			wantStatus: health.StatusUnhealthy,
			wantReason: health.ReasonFatal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			registry := healthregistry.NewBuilder()
			mustRegisterExecutionCheck(t, registry, health.TargetReady, tc.checkName, tc.fn)

			evaluator := mustExecutionEvaluator(
				t,
				registry,
				WithDefaultTimeout(0),
				WithTargetParallelChecks(health.TargetReady, 2),
			)

			report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
			if err != nil {
				t.Fatalf("Evaluate() = %v, want nil", err)
			}
			if len(report.Checks) != 1 {
				t.Fatalf("checks = %d, want 1", len(report.Checks))
			}

			res := report.Checks[0]
			if res.Status != tc.wantStatus {
				t.Fatalf("health.Status = %s, want %s", res.Status, tc.wantStatus)
			}
			if res.Reason != tc.wantReason {
				t.Fatalf("health.Reason = %s, want %s", res.Reason, tc.wantReason)
			}
			if res.Observed.IsZero() {
				t.Fatal("Observed is zero")
			}
		})
	}
}
