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
)

func TestErrInterruptedSentinel(t *testing.T) {
	if ErrInterrupted == nil {
		t.Fatalf("ErrInterrupted is nil")
	}
	if ErrInterrupted.Error() != errInterruptedMessage {
		t.Fatalf("ErrInterrupted.Error() = %q, want %q", ErrInterrupted.Error(), errInterruptedMessage)
	}
	if !errors.Is(ErrInterrupted, ErrInterrupted) {
		t.Fatalf("ErrInterrupted does not match itself")
	}
}

func TestInterrupted(t *testing.T) {
	errBoom := errors.New("boom")
	err := NewInterruptedError(errBoom)

	if !Interrupted(err) {
		t.Fatalf("Interrupted(interrupted error) = false, want true")
	}
	if !Interrupted(ErrInterrupted) {
		t.Fatalf("Interrupted(ErrInterrupted) = false, want true")
	}
	if Interrupted(errBoom) {
		t.Fatalf("Interrupted(non-interrupted error) = true, want false")
	}
	if Interrupted(context.Canceled) {
		t.Fatalf("Interrupted(context.Canceled) = true, want false")
	}
	if Interrupted(context.DeadlineExceeded) {
		t.Fatalf("Interrupted(context.DeadlineExceeded) = true, want false")
	}
	if Interrupted(nil) {
		t.Fatalf("Interrupted(nil) = true, want false")
	}
}

func TestNewInterruptedErrorWithNilCause(t *testing.T) {
	err := NewInterruptedError(nil)

	if err == nil {
		t.Fatalf("NewInterruptedError(nil) returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("NewInterruptedError(nil) does not match ErrInterrupted")
	}
	if errors.Unwrap(err) != nil {
		t.Fatalf("NewInterruptedError(nil) unwrap = %v, want nil", errors.Unwrap(err))
	}
	if err.Error() != errInterruptedMessage {
		t.Fatalf("error message = %q, want %q", err.Error(), errInterruptedMessage)
	}
}

func TestNewInterruptedErrorWithCause(t *testing.T) {
	errBoom := errors.New("boom")
	err := NewInterruptedError(errBoom)

	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("NewInterruptedError does not match ErrInterrupted")
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("NewInterruptedError does not unwrap to cause")
	}

	wantMessage := errInterruptedMessage + ": " + errBoom.Error()
	if err.Error() != wantMessage {
		t.Fatalf("error message = %q, want %q", err.Error(), wantMessage)
	}
}

func TestNewInterruptedErrorReturnsAlreadyInterruptedCause(t *testing.T) {
	inner := NewInterruptedError(errors.New("inner"))
	outer := NewInterruptedError(inner)

	if outer != inner {
		t.Fatalf("NewInterruptedError did not return already interrupted cause unchanged")
	}
}

func TestInterruptedErrorDoesNotMatchErrExhausted(t *testing.T) {
	err := NewInterruptedError(errors.New("boom"))
	if errors.Is(err, ErrExhausted) {
		t.Fatalf("interrupted error matched ErrExhausted")
	}
}
