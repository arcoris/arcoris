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
)

func TestEvaluateCheckHandlesNilChecker(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(t, emptyRegistry(t), WithClock(newStepClock(testObserved, testObserved)))

	res := evaluator.evaluateCheck(context.Background(), nil, 0)
	if res.Status != health.StatusUnknown || res.Reason != health.ReasonNotObserved {
		t.Fatalf("nil checker result = %+v", res)
	}
	if !errors.Is(res.Cause, health.ErrNilChecker) {
		t.Fatalf("nil checker cause = %v, want health.ErrNilChecker", res.Cause)
	}
}

func TestEvaluateCheckHandlesTypedNilChecker(t *testing.T) {
	t.Parallel()

	var checker *typedNilChecker
	evaluator := mustEvaluator(t, emptyRegistry(t), WithClock(newStepClock(testObserved, testObserved)))

	res := evaluator.evaluateCheck(context.Background(), checker, 0)
	if res.Status != health.StatusUnknown || res.Reason != health.ReasonNotObserved {
		t.Fatalf("typed nil checker result = %+v, want unknown not_observed", res)
	}
	if !errors.Is(res.Cause, health.ErrNilChecker) {
		t.Fatalf("typed nil checker cause = %v, want health.ErrNilChecker", res.Cause)
	}
	if !res.IsValid() {
		t.Fatalf("typed nil checker result IsValid() = false: %+v", res)
	}
}

type typedNilChecker struct{}

func (*typedNilChecker) Name() string {
	return "typed_nil"
}

func (*typedNilChecker) Check(context.Context) health.Result {
	panic("typed nil checker must not execute")
}
