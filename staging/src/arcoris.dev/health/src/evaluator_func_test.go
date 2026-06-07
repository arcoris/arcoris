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

package health

import (
	"context"
	"testing"
)

func TestEvaluatorFuncAdaptsFunction(t *testing.T) {
	t.Parallel()

	evaluator := EvaluatorFunc(func(ctx context.Context, target Target) (Report, error) {
		return Report{Target: target, Status: StatusHealthy}, nil
	})

	report, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Target != TargetReady || report.Status != StatusHealthy {
		t.Fatalf("report = %+v, want ready healthy", report)
	}
}

func TestEvaluatorFuncNilPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("Evaluate on nil EvaluatorFunc did not panic")
		}
	}()

	var evaluator EvaluatorFunc
	_, _ = evaluator.Evaluate(context.Background(), TargetReady)
}
