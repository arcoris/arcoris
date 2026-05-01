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

package wait

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// typedCause is a small typed test error used to verify errors.As traversal
// through wait error wrappers.
type typedCause struct {
	// message is the diagnostic text returned by Error.
	message string
}

// Error returns the typed test error message.
func (e typedCause) Error() string {
	return e.message
}

// wrappingError is a small test wrapper used to verify errors.Is traversal
// through non-wait error wrappers.
type wrappingError struct {
	// cause is the wrapped error exposed by Unwrap.
	cause error
}

// Error returns a diagnostic message for the wrapping test error.
func (e wrappingError) Error() string {
	return "wrapped: " + e.cause.Error()
}

// Unwrap returns the wrapped test error cause.
func (e wrappingError) Unwrap() error {
	return e.cause
}

// wrapForTest returns err wrapped in a non-wait error implementation.
func wrapForTest(err error) error {
	return wrappingError{cause: err}
}

// mustBeInterrupted fails the test unless err is a non-nil wait interruption.
func mustBeInterrupted(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("err is nil, want non-nil")
	}
	if !Interrupted(err) {
		t.Fatal("Interrupted(err) = false, want true")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatal("errors.Is(err, ErrInterrupted) = false, want true")
	}
}

// mustNotBeInterrupted fails the test if err is classified as a wait
// interruption.
func mustNotBeInterrupted(t *testing.T, err error) {
	t.Helper()

	if Interrupted(err) {
		t.Fatal("Interrupted(err) = true, want false")
	}
	if errors.Is(err, ErrInterrupted) {
		t.Fatal("errors.Is(err, ErrInterrupted) = true, want false")
	}
}

// mustBeTimedOut fails the test unless err is a non-nil wait timeout.
func mustBeTimedOut(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("err is nil, want non-nil")
	}
	if !TimedOut(err) {
		t.Fatal("TimedOut(err) = false, want true")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Fatal("errors.Is(err, ErrTimeout) = false, want true")
	}
}

// mustNotBeTimedOut fails the test if err is classified as a wait timeout.
func mustNotBeTimedOut(t *testing.T, err error) {
	t.Helper()

	if TimedOut(err) {
		t.Fatal("TimedOut(err) = true, want false")
	}
	if errors.Is(err, ErrTimeout) {
		t.Fatal("errors.Is(err, ErrTimeout) = true, want false")
	}
}

// mustMatch fails the test unless err matches target through errors.Is.
func mustMatch(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(err, %v) = false, want true", target)
	}
}

// mustNotMatch fails the test if err matches target through errors.Is.
func mustNotMatch(t *testing.T, err error, target error) {
	t.Helper()

	if errors.Is(err, target) {
		t.Fatalf("errors.Is(err, %v) = true, want false", target)
	}
}

// mustUnwrapTo fails the test unless errors.Unwrap(err) returns want.
func mustUnwrapTo(t *testing.T, err error, want error) {
	t.Helper()

	if got := errors.Unwrap(err); got != want {
		t.Fatalf("errors.Unwrap(err) = %v, want %v", got, want)
	}
}

// mustHaveMessage fails the test unless err.Error returns want.
func mustHaveMessage(t *testing.T, err error, want string) {
	t.Helper()

	if got := err.Error(); got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

// mustAsTypedCause fails the test unless err exposes want through errors.As.
func mustAsTypedCause(t *testing.T, err error, want typedCause) {
	t.Helper()

	var got typedCause
	if !errors.As(err, &got) {
		t.Fatal("errors.As(err, &got) = false, want true")
	}
	if got.message != want.message {
		t.Fatalf("errors.As extracted message %q, want %q", got.message, want.message)
	}
}

// mustEqualDuration fails the test unless got equals want.
func mustEqualDuration(t *testing.T, label string, got time.Duration, want time.Duration) {
	t.Helper()

	if got != want {
		t.Fatalf("%s = %s, want %s", label, got, want)
	}
}

// mustPanicWith fails the test unless fn panics with want.
func mustPanicWith(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		if got != want {
			t.Fatalf("panic = %s, want %s", fmt.Sprint(got), fmt.Sprint(want))
		}
	}()

	fn()
}

func mustReceiveError(t *testing.T, ch <-chan error) error {
	t.Helper()

	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second):
		t.Fatal("operation did not return before safety timeout")
		return nil
	}
}

func mustNotReceiveTimerValue(t *testing.T, ch <-chan time.Time) {
	t.Helper()

	select {
	case value := <-ch:
		t.Fatalf("received timer value %v, want no delivery", value)
	case <-time.After(10 * time.Millisecond):
	}
}
