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

package signals

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSignalErrorMatchesErrSignal(t *testing.T) {
	err := NewSignalError(testSIGTERM)

	if !errors.Is(err, ErrSignal) {
		t.Fatal("SignalError does not match ErrSignal")
	}

	var signalErr SignalError
	if !errors.As(err, &signalErr) {
		t.Fatal("SignalError was not recoverable with errors.As")
	}
	if !sameSignal(signalErr.Event.Signal, testSIGTERM) {
		t.Fatalf("signal = %v, want %v", signalErr.Event.Signal, testSIGTERM)
	}
}

func TestSignalErrorStringIncludesSignal(t *testing.T) {
	err := NewSignalError(testSIGINT)

	if !strings.Contains(err.Error(), testSIGINT.String()) {
		t.Fatalf("error %q does not contain signal %q", err.Error(), testSIGINT.String())
	}
}

func TestNewSignalErrorRejectsNilSignal(t *testing.T) {
	mustPanicWith(t, errNilSignalErrorSignal, func() {
		NewSignalError(nil)
	})
}

func TestCauseExtractsSignalEvent(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(NewSignalError(testSIGTERM))

	event, ok := Cause(ctx)
	if !ok {
		t.Fatal("Cause did not report signal event")
	}
	if !sameSignal(event.Signal, testSIGTERM) {
		t.Fatalf("signal = %v, want %v", event.Signal, testSIGTERM)
	}
}

func TestCauseReturnsFalseForNonSignalCause(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(context.Canceled)

	if event, ok := Cause(ctx); ok || event.Signal != nil {
		t.Fatalf("Cause = (%v, %v), want empty false", event, ok)
	}
}

func TestCauseRejectsNilContext(t *testing.T) {
	mustPanicWith(t, errNilCauseContext, func() {
		Cause(nil)
	})
}
