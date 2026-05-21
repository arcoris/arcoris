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


package panicassert

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Require runs call, fails the test if no panic occurs, and returns the
// recovered panic value unchanged.
func Require(t testing.TB, call func()) any {
	t.Helper()

	value, ok := capture(call)
	if !ok {
		t.Fatal("panic = nil, want non-nil")
	}

	return value
}

// RequireNone runs call and fails the test if any panic is recovered.
func RequireNone(t testing.TB, call func()) {
	t.Helper()

	value, ok := capture(call)
	if ok {
		t.Fatalf("panic = %v, want nil", value)
	}
}

// RequireMessage runs call and requires the recovered panic to format to want
// through fmt.Sprint.
//
// The helper intentionally matches formatted panic text rather than only string
// values so that tests can assert stable diagnostics from string, error, and
// other panic payloads through one API.
func RequireMessage(t testing.TB, want string, call func()) {
	t.Helper()

	got := message(Require(t, call))
	if got != want {
		t.Fatalf("panic message = %q, want %q", got, want)
	}
}

// RequireError runs call and requires the recovered panic to implement error.
// It returns the recovered error for follow-up assertions.
func RequireError(t testing.TB, call func()) error {
	t.Helper()

	value := Require(t, call)
	err, ok := asValue[error](value)
	if !ok {
		t.Fatalf("panic type = %T, want error", value)
	}

	return err
}

// RequireErrorIs runs call and requires the recovered panic to be an error that
// matches target through errors.Is.
func RequireErrorIs(t testing.TB, target error, call func()) error {
	t.Helper()

	if target == nil {
		t.Fatal("target error = nil, want non-nil")
	}

	err := RequireError(t, call)
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(panic, target) = false, want true; panic = %v, target = %v", err, target)
	}

	return err
}

// RequireValue runs call and requires the recovered panic to be assignable to T
// and deeply equal to want.
//
// reflect.DeepEqual keeps the helper safe for slices, maps, and other values
// that are valid panic payloads but not comparable with ==.
func RequireValue[T any](t testing.TB, want T, call func()) T {
	t.Helper()

	got, equal := valueMatches(Require(t, call), want)
	if !equal {
		t.Fatalf("panic value = %#v, want %#v", got, want)
	}

	return got
}

// RequireAs runs call and requires the recovered panic to be assignable to T.
// It returns the typed recovered value for follow-up assertions.
func RequireAs[T any](t testing.TB, call func()) T {
	t.Helper()

	value := Require(t, call)
	typed, ok := asValue[T](value)
	if !ok {
		var want T
		t.Fatalf("panic type = %T, want %T", value, want)
	}

	return typed
}

func capture(call func()) (value any, ok bool) {
	defer func() {
		value = recover()
		ok = value != nil
	}()

	call()
	return nil, false
}

func message(value any) string {
	return fmt.Sprint(value)
}

func asValue[T any](value any) (T, bool) {
	typed, ok := value.(T)
	return typed, ok
}

func valueMatches[T any](value any, want T) (T, bool) {
	got, ok := asValue[T](value)
	if !ok {
		var zero T
		return zero, false
	}

	return got, reflect.DeepEqual(got, want)
}
