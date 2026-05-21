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

package panicassert

import (
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

// RequireValue runs call and requires the recovered panic to be assignable to T
// and deeply equal to want.
//
// reflect.DeepEqual keeps the helper safe for slices, maps, and other values
// that are valid panic payloads but not comparable with ==.
func RequireValue[T any](t testing.TB, want T, call func()) T {
	t.Helper()

	got := RequireAs[T](t, call)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("panic value = %#v, want %#v", got, want)
	}

	return got
}

// RequireAs runs call and requires the recovered panic to be assignable to T.
// It returns the typed recovered value for follow-up assertions.
func RequireAs[T any](t testing.TB, call func()) T {
	t.Helper()

	value := Require(t, call)
	typed, ok := value.(T)
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
