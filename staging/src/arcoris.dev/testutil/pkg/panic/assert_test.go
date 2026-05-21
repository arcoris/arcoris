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
	"errors"
	"reflect"
	"testing"
)

type testPanic struct {
	Name string
	ID   int
}

func TestRequireReturnsRecoveredString(t *testing.T) {
	t.Parallel()

	got := Require(t, func() {
		panic("boom")
	})

	if got != "boom" {
		t.Fatalf("Require() = %#v, want %#v", got, "boom")
	}
}

func TestRequireReturnsRecoveredStruct(t *testing.T) {
	t.Parallel()

	want := testPanic{Name: "boom", ID: 7}
	got := Require(t, func() {
		panic(want)
	})

	if got != want {
		t.Fatalf("Require() = %#v, want %#v", got, want)
	}
}

func TestRequireNonePasses(t *testing.T) {
	t.Parallel()

	RequireNone(t, func() {})
}

func TestRequireMessageMatchesString(t *testing.T) {
	t.Parallel()

	RequireMessage(t, "boom", func() {
		panic("boom")
	})
}

func TestRequireMessageMatchesError(t *testing.T) {
	t.Parallel()

	RequireMessage(t, "boom", func() {
		panic(errors.New("boom"))
	})
}

func TestRequireMessageMatchesNumericValue(t *testing.T) {
	t.Parallel()

	RequireMessage(t, "42", func() {
		panic(42)
	})
}

func TestRequireValueMatchesString(t *testing.T) {
	t.Parallel()

	got := RequireValue(t, "boom", func() {
		panic("boom")
	})

	if got != "boom" {
		t.Fatalf("RequireValue() = %#v, want %#v", got, "boom")
	}
}

func TestRequireValueMatchesInt(t *testing.T) {
	t.Parallel()

	got := RequireValue(t, 42, func() {
		panic(42)
	})

	if got != 42 {
		t.Fatalf("RequireValue() = %#v, want %#v", got, 42)
	}
}

func TestRequireValueMatchesStruct(t *testing.T) {
	t.Parallel()

	want := testPanic{Name: "boom", ID: 7}
	got := RequireValue(t, want, func() {
		panic(want)
	})

	if got != want {
		t.Fatalf("RequireValue() = %#v, want %#v", got, want)
	}
}

func TestRequireValueMatchesSlice(t *testing.T) {
	t.Parallel()

	want := []int{1, 2, 3}
	got := RequireValue(t, want, func() {
		panic(want)
	})

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("RequireValue() = %#v, want %#v", got, want)
	}
}

func TestRequireValueMatchesMap(t *testing.T) {
	t.Parallel()

	want := map[string]int{"a": 1}
	got := RequireValue(t, want, func() {
		panic(want)
	})

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("RequireValue() = %#v, want %#v", got, want)
	}
}

func TestRequireAsReturnsTypedString(t *testing.T) {
	t.Parallel()

	got := RequireAs[string](t, func() {
		panic("boom")
	})

	if got != "boom" {
		t.Fatalf("RequireAs() = %#v, want %#v", got, "boom")
	}
}

func TestRequireAsReturnsTypedStruct(t *testing.T) {
	t.Parallel()

	want := testPanic{Name: "boom", ID: 7}
	got := RequireAs[testPanic](t, func() {
		panic(want)
	})

	if got != want {
		t.Fatalf("RequireAs() = %#v, want %#v", got, want)
	}
}

func TestCaptureWithoutPanic(t *testing.T) {
	t.Parallel()

	got, ok := capture(func() {})
	if ok {
		t.Fatal("capture() ok = true, want false")
	}
	if got != nil {
		t.Fatalf("capture() value = %#v, want nil", got)
	}
}

func TestCaptureWithPanic(t *testing.T) {
	t.Parallel()

	want := testPanic{Name: "boom", ID: 7}
	got, ok := capture(func() {
		panic(want)
	})
	if !ok {
		t.Fatal("capture() ok = false, want true")
	}
	if got != want {
		t.Fatalf("capture() value = %#v, want %#v", got, want)
	}
}

func TestMessageFormatsValues(t *testing.T) {
	t.Parallel()

	if got := message("boom"); got != "boom" {
		t.Fatalf("message(string) = %q, want %q", got, "boom")
	}
	if got := message(errors.New("boom")); got != "boom" {
		t.Fatalf("message(error) = %q, want %q", got, "boom")
	}
	if got := message(42); got != "42" {
		t.Fatalf("message(int) = %q, want %q", got, "42")
	}
}
