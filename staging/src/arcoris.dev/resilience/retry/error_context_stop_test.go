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

package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestContextStopErrorReturnsNilForActiveContext(t *testing.T) {
	err := contextStopError(context.Background())
	if err != nil {
		t.Fatalf("contextStopError(active context) = %v, want nil", err)
	}
}

func TestContextStopErrorClassifiesCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := contextStopError(ctx)
	if err == nil {
		t.Fatalf("contextStopError returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("contextStopError does not match ErrInterrupted")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("contextStopError does not preserve context.Canceled")
	}
	if Interrupted(context.Canceled) {
		t.Fatalf("raw context.Canceled classified as retry interruption")
	}
}

func TestContextStopErrorClassifiesDeadlineExceededContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	<-ctx.Done()

	err := contextStopError(ctx)
	if err == nil {
		t.Fatalf("contextStopError returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("contextStopError does not match ErrInterrupted")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("contextStopError does not preserve context.DeadlineExceeded")
	}
	if Interrupted(context.DeadlineExceeded) {
		t.Fatalf("raw context.DeadlineExceeded classified as retry interruption")
	}
}

func TestContextStopErrorPreservesCancelCause(t *testing.T) {
	cause := errors.New("custom cancel cause")

	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)

	err := contextStopError(ctx)
	if err == nil {
		t.Fatalf("contextStopError returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("contextStopError does not match ErrInterrupted")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("contextStopError does not preserve context.Canceled")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("contextStopError does not preserve custom cause")
	}
}

func TestContextStopErrorPreservesDeadlineCause(t *testing.T) {
	cause := errors.New("custom deadline cause")

	ctx, cancel := context.WithDeadlineCause(
		context.Background(),
		time.Now().Add(-time.Second),
		cause,
	)
	defer cancel()

	err := contextStopError(ctx)
	if err == nil {
		t.Fatalf("contextStopError returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("contextStopError does not match ErrInterrupted")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("contextStopError does not preserve context.DeadlineExceeded")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("contextStopError does not preserve custom deadline cause")
	}
}
