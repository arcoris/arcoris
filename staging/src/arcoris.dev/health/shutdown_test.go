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
	"testing"
)

func TestSourceChannelChecks(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})
	checker, err := NewShutdownCheck("shutdown", done)
	if err != nil {
		t.Fatalf("NewShutdownCheck() = %v, want nil", err)
	}
	if checker.Name() != "shutdown" {
		t.Fatalf("Name() = %q, want shutdown", checker.Name())
	}
	if res := checker.Check(context.Background()); res.Status != StatusHealthy {
		t.Fatalf("open channel status = %s, want healthy", res.Status)
	}

	close(done)
	result := checker.Check(context.Background())
	if result.Status != StatusUnhealthy || result.Reason != ReasonShuttingDown {
		t.Fatalf("closed channel result = %+v, want unhealthy shutting_down", result)
	}
}

func TestDrainChannelCheck(t *testing.T) {
	t.Parallel()

	draining := make(chan struct{})
	checker, err := NewDrainCheck("drain", draining)
	if err != nil {
		t.Fatalf("NewDrainCheck() = %v, want nil", err)
	}

	close(draining)
	res := checker.Check(context.Background())
	if res.Status != StatusUnhealthy || res.Reason != ReasonDraining {
		t.Fatalf("drain result = %+v, want unhealthy draining", res)
	}
}

func TestSourceContextChecks(t *testing.T) {
	t.Parallel()

	cause := errors.New("owner stop")
	source, cancel := context.WithCancelCause(context.Background())
	checker, err := NewContextShutdownCheck("shutdown", source)
	if err != nil {
		t.Fatalf("NewContextShutdownCheck() = %v, want nil", err)
	}
	if checker.Name() != "shutdown" {
		t.Fatalf("Name() = %q, want shutdown", checker.Name())
	}
	if res := checker.Check(context.Background()); res.Status != StatusHealthy {
		t.Fatalf("active source status = %s, want healthy", res.Status)
	}

	cancel(cause)
	result := checker.Check(context.Background())
	if result.Status != StatusUnhealthy || result.Reason != ReasonShuttingDown {
		t.Fatalf("canceled source result = %+v, want unhealthy shutting_down", result)
	}
	if !errors.Is(result.Cause, cause) {
		t.Fatalf("cause = %v, want %v", result.Cause, cause)
	}
}

func TestContextDrainCheck(t *testing.T) {
	t.Parallel()

	source, cancel := context.WithCancel(context.Background())
	checker, err := NewContextDrainCheck("drain", source)
	if err != nil {
		t.Fatalf("NewContextDrainCheck() = %v, want nil", err)
	}

	cancel()
	res := checker.Check(context.Background())
	if res.Status != StatusUnhealthy || res.Reason != ReasonDraining {
		t.Fatalf("drain result = %+v, want unhealthy draining", res)
	}
}

func TestSourceChecksRejectInvalidInputs(t *testing.T) {
	t.Parallel()

	if _, err := NewShutdownCheck("shutdown", nil); !errors.Is(err, ErrNilSourceChannel) {
		t.Fatalf("NewShutdownCheck(nil) = %v, want ErrNilSourceChannel", err)
	}
	if _, err := NewContextShutdownCheck("shutdown", nil); !errors.Is(err, ErrNilSourceContext) {
		t.Fatalf("NewContextShutdownCheck(nil) = %v, want ErrNilSourceContext", err)
	}
	if _, err := NewDrainCheck("bad-name", make(chan struct{})); !errors.Is(err, ErrInvalidCheckName) {
		t.Fatalf("NewDrainCheck(invalid name) = %v, want ErrInvalidCheckName", err)
	}
}

func TestMustSourceChecksPanic(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, ErrNilSourceChannel, func() {
		MustShutdownCheck("shutdown", nil)
	})
	mustPanicWith(t, ErrNilSourceChannel, func() {
		MustDrainCheck("drain", nil)
	})
	mustPanicWith(t, ErrNilSourceContext, func() {
		MustContextShutdownCheck("shutdown", nil)
	})
	mustPanicWith(t, ErrNilSourceContext, func() {
		MustContextDrainCheck("drain", nil)
	})
}

func TestMustSourceChecksReturnCheckerOnValidInput(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})
	source := context.Background()

	if MustShutdownCheck("shutdown", done).Name() != "shutdown" {
		t.Fatal("MustShutdownCheck returned wrong name")
	}
	if MustDrainCheck("drain", done).Name() != "drain" {
		t.Fatal("MustDrainCheck returned wrong name")
	}
	if MustContextShutdownCheck("context_shutdown", source).Name() != "context_shutdown" {
		t.Fatal("MustContextShutdownCheck returned wrong name")
	}
	if MustContextDrainCheck("context_drain", source).Name() != "context_drain" {
		t.Fatal("MustContextDrainCheck returned wrong name")
	}
}
