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

package value

import (
	"errors"
	"slices"
	"testing"
)

// requireNoError fails the test when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs asserts that err preserves target through errors.Is.
func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

// requireValueError asserts the structured construction error shape.
func requireValueError(
	t *testing.T,
	err error,
	target error,
	path string,
	reason ErrorReason,
) *Error {
	t.Helper()
	requireErrorIs(t, err, target)

	var valueErr *Error
	if !errors.As(err, &valueErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	if valueErr.Path != path {
		t.Fatalf("Error.Path = %q, want %q", valueErr.Path, path)
	}
	if valueErr.Reason != reason {
		t.Fatalf("Error.Reason = %q, want %q", valueErr.Reason, reason)
	}
	if valueErr.Detail == "" {
		t.Fatal("Error.Detail is empty")
	}

	return valueErr
}

// requireEqual compares values without hiding actual output.
func requireEqual[T comparable](t *testing.T, got T, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

// requireBytesEqual compares byte slices with clear failure output.
func requireBytesEqual(t *testing.T, got []byte, want []byte) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

// requireStringsEqual compares string slices while preserving useful output.
func requireStringsEqual(t *testing.T, got []string, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

// requirePanic asserts that fn panics.
func requirePanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	fn()
}

// mustObject constructs an object value for tests.
func mustObject(t *testing.T, fields ...Field) Value {
	t.Helper()

	value, err := NewObject(fields...)
	requireNoError(t, err)

	return value
}

// mustList constructs a list value for tests.
func mustList(t *testing.T, items ...Value) Value {
	t.Helper()

	value, err := NewList(items...)
	requireNoError(t, err)

	return value
}

// mustMap constructs a map value for tests.
func mustMap(t *testing.T, entries ...Entry) Value {
	t.Helper()

	value, err := NewMap(entries...)
	requireNoError(t, err)

	return value
}
