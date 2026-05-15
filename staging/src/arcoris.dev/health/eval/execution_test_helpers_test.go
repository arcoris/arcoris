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

package eval

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"arcoris.dev/health"
)

const executionTestTimeout = 5 * time.Second

func mustRegisterExecutionCheck(
	t *testing.T,
	r *health.Registry,
	target health.Target,
	name string,
	fn health.CheckFunc,
) {
	t.Helper()

	chk, err := health.NewCheck(name, fn)
	if err != nil {
		t.Fatalf("health.NewCheck(%q) = %v, want nil", name, err)
	}

	if err := r.Register(target, chk); err != nil {
		t.Fatalf("Register(%s, %q) = %v, want nil", target, name, err)
	}
}

func mustExecutionEvaluator(
	t *testing.T,
	r *health.Registry,
	opts ...EvaluatorOption,
) *Evaluator {
	t.Helper()

	evaluator, err := NewEvaluator(r, opts...)
	if err != nil {
		t.Fatalf("NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}

func executionResultNames(results []health.Result) []string {
	names := make([]string, 0, len(results))
	for _, res := range results {
		names = append(names, res.Name)
	}

	return names
}

func executionCheckName(i int) string {
	names := []string{
		"check_zero",
		"check_one",
		"check_two",
		"check_three",
		"check_four",
		"check_five",
		"check_six",
		"check_seven",
		"check_eight",
		"check_nine",
	}
	if i >= 0 && i < len(names) {
		return names[i]
	}

	return "check_extra"
}

func updateMaxInt64(maxSeen *atomic.Int64, cur int64) {
	for {
		prev := maxSeen.Load()
		if cur <= prev {
			return
		}
		if maxSeen.CompareAndSwap(prev, cur) {
			return
		}
	}
}

func sameStrings(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

func blockingAfterContextDone(release <-chan struct{}) health.CheckFunc {
	return func(ctx context.Context) health.Result {
		<-ctx.Done()
		<-release
		return health.Healthy("blocking_check")
	}
}
