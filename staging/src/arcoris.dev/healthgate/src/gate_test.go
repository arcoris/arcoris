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

package healthgate

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestNewValidatesAndNormalizesInitialResult(t *testing.T) {
	t.Parallel()

	gate, err := New("ready_gate", health.Result{Status: health.StatusHealthy})
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

	result := gate.Check(context.Background())
	if result.Name != "ready_gate" || result.Status != health.StatusHealthy {
		t.Fatalf("initial result = %+v, want named healthy result", result)
	}
}

func TestNewRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func() error
		want error
	}{
		{
			name: "invalid gate name",
			call: func() error {
				_, err := New("bad-name", health.Healthy("bad-name"))
				return err
			},
			want: health.ErrInvalidCheckName,
		},
		{
			name: "invalid result",
			call: func() error {
				_, err := New("ready_gate", health.Result{Status: health.Status(99)})
				return err
			},
			want: ErrInvalidGateResult,
		},
		{
			name: "invalid result name",
			call: func() error {
				_, err := New("ready_gate", health.Healthy("bad-name"))
				return err
			},
			want: ErrMismatchedGateResult,
		},
		{
			name: "mismatched result name",
			call: func() error {
				_, err := New("ready_gate", health.Healthy("other_gate"))
				return err
			},
			want: ErrMismatchedGateResult,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tc.call(); !errors.Is(err, tc.want) {
				t.Fatalf("New() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestNewUnknown(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknown("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknown() = %v, want nil", err)
	}

	result := gate.Snapshot()
	if result.Name != "ready_gate" ||
		result.Status != health.StatusUnknown ||
		result.Reason != health.ReasonNotObserved {
		t.Fatalf("Snapshot() = %+v, want initial unknown", result)
	}
}

func TestSetRejectsInvalidResultWithoutMutation(t *testing.T) {
	t.Parallel()

	gate := mustGate(t)

	if err := gate.Set(health.Result{Status: health.Status(99)}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("Set(invalid) = %v, want ErrInvalidGateResult", err)
	}
	if got := gate.Snapshot().Status; got != health.StatusHealthy {
		t.Fatalf("status after invalid Set = %s, want healthy", got)
	}

	if err := gate.Set(health.Healthy("other_gate")); !errors.Is(err, ErrMismatchedGateResult) {
		t.Fatalf("Set(mismatch) = %v, want ErrMismatchedGateResult", err)
	}
	if got := gate.Snapshot().Name; got != "ready_gate" {
		t.Fatalf("name after mismatched Set = %q, want ready_gate", got)
	}
}

func TestNilGateBehavior(t *testing.T) {
	t.Parallel()

	var gate *Gate
	mutations := []func() error{
		func() error { return gate.Set(health.Healthy("ready_gate")) },
		gate.Healthy,
		func() error { return gate.Unknown(health.ReasonNotObserved, "unknown") },
		func() error { return gate.Starting(health.ReasonStarting, "starting") },
		func() error { return gate.Degraded(health.ReasonOverloaded, "overloaded") },
		func() error { return gate.Unhealthy(health.ReasonFatal, "fatal") },
	}

	for index, mutate := range mutations {
		if err := mutate(); !errors.Is(err, health.ErrNilChecker) {
			t.Fatalf("mutation %d = %v, want ErrNilChecker", index, err)
		}
	}

	result := gate.Check(context.Background())
	if result.Status != health.StatusUnknown || !errors.Is(result.Cause, health.ErrNilChecker) {
		t.Fatalf("nil Check() = %+v, want unknown ErrNilChecker", result)
	}
}

func TestStatusHelpersPublishExpectedResults(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknown("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknown() = %v, want nil", err)
	}

	tests := []struct {
		name   string
		set    func() error
		status health.Status
		reason health.Reason
	}{
		{"healthy", gate.Healthy, health.StatusHealthy, health.ReasonNone},
		{"unknown", func() error { return gate.Unknown(health.ReasonNotObserved, "unknown") }, health.StatusUnknown, health.ReasonNotObserved},
		{"starting", func() error { return gate.Starting(health.ReasonStarting, "starting") }, health.StatusStarting, health.ReasonStarting},
		{"degraded", func() error { return gate.Degraded(health.ReasonOverloaded, "degraded") }, health.StatusDegraded, health.ReasonOverloaded},
		{"unhealthy", func() error { return gate.Unhealthy(health.ReasonFatal, "fatal") }, health.StatusUnhealthy, health.ReasonFatal},
	}

	for _, tc := range tests {
		if err := tc.set(); err != nil {
			t.Fatalf("%s set = %v, want nil", tc.name, err)
		}

		result := gate.Check(context.Background())
		if result.Name != "ready_gate" || result.Status != tc.status || result.Reason != tc.reason {
			t.Fatalf("%s result = %+v, want status %s reason %s", tc.name, result, tc.status, tc.reason)
		}
	}
}

func TestConcurrentCheckAndSetRaceFree(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknown("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknown() = %v, want nil", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if i%2 == 0 {
				_ = gate.Healthy()
			} else {
				_ = gate.Degraded(health.ReasonOverloaded, "overloaded")
			}
			_ = gate.Check(context.Background())
			_ = gate.Snapshot()
		}(i)
	}
	wg.Wait()
}

func TestInvalidGateResultErrorClassification(t *testing.T) {
	t.Parallel()

	invalid := InvalidGateResultError{
		GateName: "ready_gate",
		Result:   health.Result{Status: health.Status(99), Duration: -time.Second},
	}
	if !errors.Is(invalid, ErrInvalidGateResult) {
		t.Fatal("InvalidGateResultError should match ErrInvalidGateResult")
	}

	mismatch := MismatchedGateResultError{
		GateName:   "ready_gate",
		ResultName: "other_gate",
	}
	if !errors.Is(mismatch, ErrMismatchedGateResult) {
		t.Fatal("MismatchedGateResultError should match ErrMismatchedGateResult")
	}
}

func mustGate(t *testing.T) *Gate {
	t.Helper()

	gate, err := New("ready_gate", health.Healthy("ready_gate"))
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

	return gate
}
