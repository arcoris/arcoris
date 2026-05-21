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

package errorassert

import (
	"errors"
	"testing"
)

// Require fails the test unless err is non-nil.
func Require(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("err = nil, want non-nil")
	}
}

// RequireNone fails the test if err is non-nil.
func RequireNone(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

// RequireIs fails the test unless errors.Is(err, target) is true.
func RequireIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(err, target) = false, want true; err = %v, target = %v", err, target)
	}
}

// RequireIsNot fails the test if errors.Is(err, target) is true.
func RequireIsNot(t testing.TB, err error, target error) {
	t.Helper()

	if errors.Is(err, target) {
		t.Fatalf("errors.Is(err, target) = true, want false; err = %v, target = %v", err, target)
	}
}

// RequireMessage fails the test unless err is non-nil and err.Error() equals
// want.
func RequireMessage(t testing.TB, err error, want string) {
	t.Helper()

	Require(t, err)
	if got := err.Error(); got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

// RequireUnwrapsTo fails the test unless errors.Unwrap(err) returns want
// directly. It intentionally does not use errors.Is.
func RequireUnwrapsTo(t testing.TB, err error, want error) {
	t.Helper()

	if got := errors.Unwrap(err); got != want {
		t.Fatalf("errors.Unwrap(err) = %v, want %v", got, want)
	}
}

// RequireAs fails the test unless errors.As can extract T from err. It returns
// the extracted value for follow-up assertions.
func RequireAs[T any](t testing.TB, err error) T {
	t.Helper()

	var got T
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, %T) = false, want true; err = %v", got, err)
	}

	return got
}
