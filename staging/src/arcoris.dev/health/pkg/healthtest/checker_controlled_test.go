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
	"sync"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestBlockingCheckerRelease(t *testing.T) {
	t.Parallel()

	checker := NewBlockingChecker("storage")
	done := make(chan health.Result, 1)
	go func() {
		done <- checker.Check(context.Background())
	}()

	mustClose(t, checker.Started())
	checker.Release(HealthyResult("storage"))

	select {
	case res := <-done:
		AssertResultStatus(t, res, health.StatusHealthy)
	case <-time.After(time.Second):
		t.Fatal("blocking checker did not return")
	}
	if checker.Calls() != 1 {
		t.Fatalf("Calls = %d, want 1", checker.Calls())
	}
}

func TestBlockingCheckerCancellation(t *testing.T) {
	t.Parallel()

	checker := NewBlockingChecker("storage")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan health.Result, 1)
	go func() {
		done <- checker.Check(ctx)
	}()

	mustClose(t, checker.Started())
	cancel()

	select {
	case res := <-done:
		AssertResultStatus(t, res, health.StatusUnknown)
		AssertResultReason(t, res, health.ReasonCanceled)
	case <-time.After(time.Second):
		t.Fatal("blocking checker did not observe cancellation")
	}
}

func TestBlockingCheckerConcurrentCalls(t *testing.T) {
	t.Parallel()

	checker := NewBlockingChecker("storage")
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			AssertResultStatus(t, checker.Check(context.Background()), health.StatusHealthy)
		}()
	}

	mustClose(t, checker.Started())
	checker.Release(HealthyResult("storage"))
	wg.Wait()
	if checker.Calls() != 8 {
		t.Fatalf("Calls = %d, want 8", checker.Calls())
	}
}

func TestSequenceChecker(t *testing.T) {
	t.Parallel()

	checker := NewSequenceChecker(
		"storage",
		HealthyResult("storage"),
		DegradedResult("storage", health.ReasonOverloaded),
	)

	AssertResultStatus(t, checker.Check(context.Background()), health.StatusHealthy)
	AssertResultStatus(t, checker.Check(context.Background()), health.StatusDegraded)
	AssertResultStatus(t, checker.Check(context.Background()), health.StatusDegraded)
	if checker.Calls() != 3 {
		t.Fatalf("Calls = %d, want 3", checker.Calls())
	}
}

func TestSequenceCheckerFillsEmptyResultName(t *testing.T) {
	t.Parallel()

	checker := NewSequenceChecker("storage", health.Result{Status: health.StatusHealthy})
	res := checker.Check(context.Background())
	if res.Name != "storage" {
		t.Fatalf("result name = %q, want storage", res.Name)
	}
}

func mustClose(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("channel did not close")
	}
}
