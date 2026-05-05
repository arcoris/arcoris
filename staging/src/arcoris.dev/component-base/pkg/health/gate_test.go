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
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestGateConstructionAndCheck(t *testing.T) {
	t.Parallel()

	gate, err := NewGate("ready_gate", Result{Status: StatusHealthy})
	if err != nil {
		t.Fatalf("NewGate() = %v, want nil", err)
	}
	if got := gate.Name(); got != "ready_gate" {
		t.Fatalf("Name() = %q, want ready_gate", got)
	}
	if result := gate.Check(context.Background()); result.Name != "ready_gate" || result.Status != StatusHealthy {
		t.Fatalf("Check() = %+v, want named healthy", result)
	}
}

func TestGateSetters(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknownGate("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknownGate() = %v, want nil", err)
	}

	setters := []struct {
		name   string
		set    func() error
		status Status
	}{
		{"unknown", func() error { return gate.Unknown(ReasonNotObserved, "unknown") }, StatusUnknown},
		{"starting", func() error { return gate.Starting(ReasonStarting, "starting") }, StatusStarting},
		{"healthy", gate.Healthy, StatusHealthy},
		{"degraded", func() error { return gate.Degraded(ReasonOverloaded, "overloaded") }, StatusDegraded},
		{"unhealthy", func() error { return gate.Unhealthy(ReasonFatal, "fatal") }, StatusUnhealthy},
	}

	for _, tc := range setters {
		if err := tc.set(); err != nil {
			t.Fatalf("%s setter = %v, want nil", tc.name, err)
		}
		if got := gate.Check(context.Background()).Status; got != tc.status {
			t.Fatalf("%s status = %s, want %s", tc.name, got, tc.status)
		}
	}
}

func TestGateRejectsInvalidResults(t *testing.T) {
	t.Parallel()

	if _, err := NewGate("bad-name", Healthy("bad-name")); !errors.Is(err, ErrInvalidCheckName) {
		t.Fatalf("NewGate(invalid name) = %v, want ErrInvalidCheckName", err)
	}
	if _, err := NewGate("ready_gate", Result{Status: Status(99)}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("NewGate(invalid result) = %v, want ErrInvalidGateResult", err)
	}
	if _, err := NewGate("ready_gate", Result{Status: StatusHealthy, Reason: Reason("bad-reason")}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("NewGate(invalid reason) = %v, want ErrInvalidGateResult", err)
	}
	if _, err := NewGate("ready_gate", Healthy("other_gate")); !errors.Is(err, ErrMismatchedGateResult) {
		t.Fatalf("NewGate(mismatched result) = %v, want ErrMismatchedGateResult", err)
	}

	gate, err := NewUnknownGate("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknownGate() = %v, want nil", err)
	}
	if err := gate.Set(Result{Status: StatusHealthy, Duration: -time.Second}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("Set(invalid result) = %v, want ErrInvalidGateResult", err)
	}
	if err := gate.Set(Result{Status: StatusHealthy, Reason: Reason("bad-reason")}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("Set(invalid reason) = %v, want ErrInvalidGateResult", err)
	}
	if err := gate.Set(Healthy("other_gate")); !errors.Is(err, ErrMismatchedGateResult) {
		t.Fatalf("Set(mismatch) = %v, want ErrMismatchedGateResult", err)
	}
}

func TestNilGateBehavior(t *testing.T) {
	t.Parallel()

	var gate *Gate
	if gate.Name() != "" {
		t.Fatal("nil Gate Name() should be empty")
	}
	if result := gate.Check(context.Background()); !errors.Is(result.Cause, ErrNilChecker) {
		t.Fatalf("nil Gate Check cause = %v, want ErrNilChecker", result.Cause)
	}
	if err := gate.Set(Healthy("ready_gate")); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("nil Gate Set() = %v, want ErrNilChecker", err)
	}
	if err := gate.Healthy(); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("nil Gate Healthy() = %v, want ErrNilChecker", err)
	}
	if err := gate.Unknown(ReasonNotObserved, "unknown"); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("nil Gate Unknown() = %v, want ErrNilChecker", err)
	}
	if err := gate.Starting(ReasonStarting, "starting"); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("nil Gate Starting() = %v, want ErrNilChecker", err)
	}
	if err := gate.Degraded(ReasonOverloaded, "overloaded"); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("nil Gate Degraded() = %v, want ErrNilChecker", err)
	}
	if err := gate.Unhealthy(ReasonFatal, "fatal"); !errors.Is(err, ErrNilChecker) {
		t.Fatalf("nil Gate Unhealthy() = %v, want ErrNilChecker", err)
	}
}

func TestGateConcurrentAccess(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknownGate("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknownGate() = %v, want nil", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := gate.Healthy(); err != nil {
				t.Errorf("Healthy() = %v, want nil", err)
			}
			_ = gate.Check(context.Background())
		}()
	}
	wg.Wait()
}

func TestGateErrors(t *testing.T) {
	t.Parallel()

	invalid := InvalidGateResultError{GateName: "gate", Result: Result{Status: Status(99)}}
	mismatch := MismatchedGateResultError{GateName: "gate", ResultName: "other"}

	mustErrorIs(t, invalid, ErrInvalidGateResult)
	mustErrorIs(t, mismatch, ErrMismatchedGateResult)
	if invalid.Error() == "" || mismatch.Error() == "" {
		t.Fatal("gate error messages must be non-empty")
	}
}
