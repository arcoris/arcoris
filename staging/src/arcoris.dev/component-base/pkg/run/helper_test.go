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

package run

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestWaitReturnsContextCause(t *testing.T) {
	t.Parallel()

	want := errors.New("stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	if err := Wait(ctx); !errors.Is(err, want) {
		t.Fatalf("Wait error = %v, want %v", err, want)
	}
}

func TestWaitFallsBackToContextErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := Wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait error = %v, want context.Canceled", err)
	}
}

func TestWaitRejectsNilContext(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilWaitContext, func() {
		Wait(nil)
	})
}

func TestIsContextStop(t *testing.T) {
	t.Parallel()

	if !IsContextStop(context.Canceled) {
		t.Fatal("context.Canceled was not classified as context stop")
	}
	if !IsContextStop(context.DeadlineExceeded) {
		t.Fatal("context.DeadlineExceeded was not classified as context stop")
	}
	if IsContextStop(errors.New("other")) {
		t.Fatal("unrelated error was classified as context stop")
	}
}

func TestIgnoreContextCanceled(t *testing.T) {
	t.Parallel()

	if err := IgnoreContextCanceled(context.Canceled); err != nil {
		t.Fatalf("IgnoreContextCanceled = %v, want nil", err)
	}
	if err := IgnoreContextCanceled(context.DeadlineExceeded); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("IgnoreContextCanceled = %v, want deadline", err)
	}

	errOther := errors.New("other")
	if err := IgnoreContextCanceled(errOther); !errors.Is(err, errOther) {
		t.Fatalf("IgnoreContextCanceled = %v, want other", err)
	}
}

func TestIgnoreContextStopIgnoresObservedContextCause(t *testing.T) {
	t.Parallel()

	want := errors.New("owner stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	if err := IgnoreContextStop(ctx, want); err != nil {
		t.Fatalf("IgnoreContextStop = %v, want nil", err)
	}
}

func TestIgnoreContextStopIgnoresContextErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := IgnoreContextStop(ctx, ctx.Err()); err != nil {
		t.Fatalf("IgnoreContextStop = %v, want nil", err)
	}
}

func TestIgnoreContextStopIgnoresCanceledContextErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := IgnoreContextStop(ctx, context.Canceled); err != nil {
		t.Fatalf("IgnoreContextStop(context.Canceled) = %v, want nil", err)
	}
}

func TestIgnoreContextStopIgnoresDeadlineContextErr(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	if err := IgnoreContextStop(ctx, context.DeadlineExceeded); err != nil {
		t.Fatalf("IgnoreContextStop(context.DeadlineExceeded) = %v, want nil", err)
	}
}

func TestIgnoreContextStopPreservesDeadlineForCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := IgnoreContextStop(ctx, context.DeadlineExceeded)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("IgnoreContextStop = %v, want context.DeadlineExceeded", err)
	}
}

func TestIgnoreContextStopPreservesCanceledForDeadlineContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	err := IgnoreContextStop(ctx, context.Canceled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("IgnoreContextStop = %v, want context.Canceled", err)
	}
}

func TestIgnoreContextStopIgnoresWrappedContextCause(t *testing.T) {
	t.Parallel()

	want := errors.New("owner stop")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(want)

	if err := IgnoreContextStop(ctx, fmt.Errorf("wrapped: %w", want)); err != nil {
		t.Fatalf("IgnoreContextStop = %v, want nil", err)
	}
}

func TestIgnoreContextStopPreservesErrorsBeforeContextStops(t *testing.T) {
	t.Parallel()

	err := IgnoreContextStop(context.Background(), context.Canceled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("IgnoreContextStop = %v, want context.Canceled", err)
	}
}

func TestIgnoreContextStopPreservesUnrelatedErrors(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("owner stop"))

	errOther := errors.New("other")
	if err := IgnoreContextStop(ctx, errOther); !errors.Is(err, errOther) {
		t.Fatalf("IgnoreContextStop = %v, want other", err)
	}
}

func TestIgnoreContextStopRejectsNilContext(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilIgnoreContext, func() {
		IgnoreContextStop(nil, context.Canceled)
	})
}
