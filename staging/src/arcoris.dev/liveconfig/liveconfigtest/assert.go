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

package liveconfigtest

import "reflect"

// TestingT is the minimal testing interface used by assertion helpers.
type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

// RequireValue fails the test when got and want are not equal.
//
// When equal is nil, RequireValue uses reflect.DeepEqual. Tests that care about
// domain-specific equality should pass an explicit equality function.
func RequireValue[T any](t TestingT, got, want T, equal func(T, T) bool) {
	t.Helper()

	if equal == nil {
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("value = %#v, want %#v", got, want)
		}
		return
	}

	if !equal(got, want) {
		t.Fatalf("value = %#v, want %#v", got, want)
	}
}

// RequireConfigEqual fails the test when got and want differ according to
// EqualConfig.
func RequireConfigEqual(t TestingT, got, want Config) {
	t.Helper()
	RequireValue(t, got, want, EqualConfig)
}
