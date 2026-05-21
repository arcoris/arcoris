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
	"fmt"
	"reflect"
	"testing"
)

type testPanic struct {
	Name   string
	ID     int
	Nested struct {
		Flags []string
	}
}

type stringerPanic struct {
	Code int
}

func (p stringerPanic) String() string {
	return fmt.Sprintf("stringer:%d", p.Code)
}

func TestRequireReturnsRecoveredValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want any
	}{
		{
			name: "string",
			want: "boom",
		},
		{
			name: "struct",
			want: func() testPanic {
				value := testPanic{Name: "boom", ID: 7}
				value.Nested.Flags = []string{"x", "y"}
				return value
			}(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Require(t, func() {
				panic(tt.want)
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Require() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRequireNonePasses(t *testing.T) {
	t.Parallel()

	RequireNone(t, func() {})
}

func TestRequireMessageMatchesFormattedValue(t *testing.T) {
	t.Parallel()

	boom := errors.New("boom")

	tests := []struct {
		name  string
		want  string
		value any
	}{
		{name: "string", want: "boom", value: "boom"},
		{name: "error", want: "boom", value: boom},
		{name: "int", want: "42", value: 42},
		{name: "stringer", want: "stringer:7", value: stringerPanic{Code: 7}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			RequireMessage(t, tt.want, func() {
				panic(tt.value)
			})
		})
	}
}

func TestRequireErrorReturnsRecoveredError(t *testing.T) {
	t.Parallel()

	want := errors.New("boom")
	got := RequireError(t, func() {
		panic(want)
	})

	if got != want {
		t.Fatalf("RequireError() = %v, want %v", got, want)
	}
}

func TestRequireErrorIs(t *testing.T) {
	t.Parallel()

	t.Run("sentinel", func(t *testing.T) {
		t.Parallel()

		want := errors.New("sentinel")
		got := RequireErrorIs(t, want, func() {
			panic(want)
		})

		if got != want {
			t.Fatalf("RequireErrorIs() = %v, want %v", got, want)
		}
	})

	t.Run("wrapped", func(t *testing.T) {
		t.Parallel()

		want := errors.New("sentinel")
		wrapped := fmt.Errorf("wrapped: %w", want)
		got := RequireErrorIs(t, want, func() {
			panic(wrapped)
		})

		if got != wrapped {
			t.Fatalf("RequireErrorIs() = %v, want %v", got, wrapped)
		}
	})
}

func TestRequireValueMatches(t *testing.T) {
	t.Parallel()

	typedNil := (*testPanic)(nil)
	structValue := testPanic{Name: "boom", ID: 7}
	structValue.Nested.Flags = []string{"x", "y"}
	structPointer := &structValue
	testErr := errors.New("boom")

	tests := []struct {
		name string
		want any
	}{
		{name: "string", want: "boom"},
		{name: "int", want: 42},
		{name: "pointer", want: structPointer},
		{name: "typed nil pointer", want: typedNil},
		{name: "slice", want: []int{1, 2, 3}},
		{name: "map", want: map[string]int{"a": 1, "b": 2}},
		{name: "nested struct", want: structValue},
		{name: "error", want: testErr},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch want := tt.want.(type) {
			case string:
				got := RequireValue(t, want, func() { panic(want) })
				if got != want {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			case int:
				got := RequireValue(t, want, func() { panic(want) })
				if got != want {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			case *testPanic:
				got := RequireValue(t, want, func() { panic(want) })
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			case []int:
				got := RequireValue(t, want, func() { panic(want) })
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			case map[string]int:
				got := RequireValue(t, want, func() { panic(want) })
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			case testPanic:
				got := RequireValue(t, want, func() { panic(want) })
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			case error:
				got := RequireValue(t, want, func() { panic(want) })
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("RequireValue() = %#v, want %#v", got, want)
				}
			default:
				t.Fatalf("unsupported test type %T", want)
			}
		})
	}
}

func TestRequireAsReturnsTypedValue(t *testing.T) {
	t.Parallel()

	typedErr := errors.New("boom")
	structValue := testPanic{Name: "boom", ID: 7}
	structPointer := &structValue

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		got := RequireAs[string](t, func() {
			panic("boom")
		})

		if got != "boom" {
			t.Fatalf("RequireAs() = %#v, want %#v", got, "boom")
		}
	})

	t.Run("struct", func(t *testing.T) {
		t.Parallel()

		got := RequireAs[testPanic](t, func() {
			panic(structValue)
		})

		if !reflect.DeepEqual(got, structValue) {
			t.Fatalf("RequireAs() = %#v, want %#v", got, structValue)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		t.Parallel()

		got := RequireAs[*testPanic](t, func() {
			panic(structPointer)
		})

		if got != structPointer {
			t.Fatalf("RequireAs() = %#v, want %#v", got, structPointer)
		}
	})

	t.Run("error interface", func(t *testing.T) {
		t.Parallel()

		got := RequireAs[error](t, func() {
			panic(typedErr)
		})

		if !reflect.DeepEqual(got, typedErr) {
			t.Fatalf("RequireAs() = %#v, want %#v", got, typedErr)
		}
	})

	t.Run("empty interface", func(t *testing.T) {
		t.Parallel()

		got := RequireAs[any](t, func() {
			panic(structPointer)
		})

		if got != structPointer {
			t.Fatalf("RequireAs() = %#v, want %#v", got, structPointer)
		}
	})
}

func TestCapture(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")
	sliceValue := []int{1, 2, 3}
	mapValue := map[string]int{"a": 1}

	tests := []struct {
		name   string
		call   func()
		want   any
		wantOK bool
	}{
		{
			name:   "no panic",
			call:   func() {},
			want:   nil,
			wantOK: false,
		},
		{
			name:   "string panic",
			call:   func() { panic("boom") },
			want:   "boom",
			wantOK: true,
		},
		{
			name:   "struct panic",
			call:   func() { panic(testPanic{Name: "boom", ID: 7}) },
			want:   testPanic{Name: "boom", ID: 7},
			wantOK: true,
		},
		{
			name:   "error panic",
			call:   func() { panic(errBoom) },
			want:   errBoom,
			wantOK: true,
		},
		{
			name:   "slice panic",
			call:   func() { panic(sliceValue) },
			want:   sliceValue,
			wantOK: true,
		},
		{
			name:   "map panic",
			call:   func() { panic(mapValue) },
			want:   mapValue,
			wantOK: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := capture(tt.call)
			if ok != tt.wantOK {
				t.Fatalf("capture() ok = %v, want %v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("capture() value = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "nil", value: nil, want: "<nil>"},
		{name: "string", value: "boom", want: "boom"},
		{name: "error", value: errors.New("boom"), want: "boom"},
		{name: "int", value: 42, want: "42"},
		{name: "stringer", value: stringerPanic{Code: 9}, want: "stringer:9"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := message(tt.value); got != tt.want {
				t.Fatalf("message() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAsValue(t *testing.T) {
	t.Parallel()

	typedErr := errors.New("boom")
	structValue := testPanic{Name: "boom", ID: 7}
	structPointer := &structValue
	typedNil := (*testPanic)(nil)

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "matching concrete type",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := asValue[testPanic](structValue)
				if !ok || !reflect.DeepEqual(got, structValue) {
					t.Fatalf("asValue() = (%#v, %v), want (%#v, true)", got, ok, structValue)
				}
			},
		},
		{
			name: "matching pointer type",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := asValue[*testPanic](structPointer)
				if !ok || got != structPointer {
					t.Fatalf("asValue() = (%#v, %v), want (%#v, true)", got, ok, structPointer)
				}
			},
		},
		{
			name: "typed nil pointer",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := asValue[*testPanic](typedNil)
				if !ok || got != nil {
					t.Fatalf("asValue() = (%#v, %v), want (nil, true)", got, ok)
				}
			},
		},
		{
			name: "matching interface",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := asValue[error](typedErr)
				if !ok || !reflect.DeepEqual(got, typedErr) {
					t.Fatalf("asValue() = (%#v, %v), want (%#v, true)", got, ok, typedErr)
				}
			},
		},
		{
			name: "mismatched type",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := asValue[int]("boom")
				if ok || got != 0 {
					t.Fatalf("asValue() = (%#v, %v), want (0, false)", got, ok)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t)
		})
	}
}

func TestValueMatches(t *testing.T) {
	t.Parallel()

	typedNil := (*testPanic)(nil)

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "equal slices",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := valueMatches([]int{1, 2}, []int{1, 2})
				if !ok || !reflect.DeepEqual(got, []int{1, 2}) {
					t.Fatalf("valueMatches() = (%#v, %v), want (%#v, true)", got, ok, []int{1, 2})
				}
			},
		},
		{
			name: "different slices",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := valueMatches([]int{1, 2}, []int{1, 3})
				if ok || !reflect.DeepEqual(got, []int{1, 2}) {
					t.Fatalf("valueMatches() = (%#v, %v), want (%#v, false)", got, ok, []int{1, 2})
				}
			},
		},
		{
			name: "equal maps",
			check: func(t *testing.T) {
				t.Helper()
				want := map[string]int{"a": 1}
				got, ok := valueMatches(want, map[string]int{"a": 1})
				if !ok || !reflect.DeepEqual(got, want) {
					t.Fatalf("valueMatches() = (%#v, %v), want (%#v, true)", got, ok, want)
				}
			},
		},
		{
			name: "different maps",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := valueMatches(map[string]int{"a": 1}, map[string]int{"a": 2})
				if ok || !reflect.DeepEqual(got, map[string]int{"a": 1}) {
					t.Fatalf("valueMatches() = (%#v, %v), want (%#v, false)", got, ok, map[string]int{"a": 1})
				}
			},
		},
		{
			name: "typed nil pointer equality",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := valueMatches[*testPanic](typedNil, typedNil)
				if !ok || got != nil {
					t.Fatalf("valueMatches() = (%#v, %v), want (nil, true)", got, ok)
				}
			},
		},
		{
			name: "type mismatch",
			check: func(t *testing.T) {
				t.Helper()
				got, ok := valueMatches[int]("boom", 42)
				if ok || got != 0 {
					t.Fatalf("valueMatches() = (%#v, %v), want (0, false)", got, ok)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t)
		})
	}
}
