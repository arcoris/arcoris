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
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/healthregistry"
)

func TestEvaluatorParallelTimeoutAndCancel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ctx        context.Context
		timeout    time.Duration
		wantReason health.Reason
	}{
		{
			name:       "timeout",
			ctx:        context.Background(),
			timeout:    time.Nanosecond,
			wantReason: health.ReasonTimeout,
		},
		{
			name: "canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			timeout:    time.Second,
			wantReason: health.ReasonCanceled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			release := make(chan struct{})
			defer close(release)

			registry := healthregistry.NewBuilder()
			mustRegisterExecutionCheck(t, registry, health.TargetReady, "blocking_one", blockingAfterContextDone(release))
			mustRegisterExecutionCheck(t, registry, health.TargetReady, "blocking_two", blockingAfterContextDone(release))

			evaluator := mustExecutionEvaluator(
				t,
				registry,
				WithTargetTimeout(health.TargetReady, tc.timeout),
				WithTargetParallelChecks(health.TargetReady, 2),
			)

			done := make(chan health.Report, 1)
			go func() {
				report, err := evaluator.Evaluate(tc.ctx, health.TargetReady)
				if err != nil {
					t.Errorf("Evaluate() = %v, want nil", err)
				}
				done <- report
			}()

			var report health.Report
			select {
			case report = <-done:
			case <-time.After(executionTestTimeout):
				t.Fatal("parallel evaluation did not finish")
			}

			if len(report.Checks) != 2 {
				t.Fatalf("checks = %d, want 2", len(report.Checks))
			}
			for _, res := range report.Checks {
				if res.Reason != tc.wantReason {
					t.Fatalf("health.Reason for %s = %s, want %s", res.Name, res.Reason, tc.wantReason)
				}
			}
		})
	}
}
