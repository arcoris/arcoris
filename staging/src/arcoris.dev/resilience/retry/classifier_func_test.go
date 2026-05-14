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
	"errors"
	"testing"
)

func TestClassifierFuncPanicsWhenNil(t *testing.T) {
	var classifier ClassifierFunc

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("ClassifierFunc.Retryable did not panic")
		}
		if recovered != panicNilClassifierFunc {
			t.Fatalf("panic = %v, want %q", recovered, panicNilClassifierFunc)
		}
	}()

	_ = classifier.Retryable(errors.New("boom"))
}

func TestClassifierFuncRejectsNilWithoutCallingWrappedFunction(t *testing.T) {
	called := false

	classifier := ClassifierFunc(func(error) bool {
		called = true
		return true
	})

	if classifier.Retryable(nil) {
		t.Fatalf("ClassifierFunc.Retryable(nil) = true, want false")
	}
	if called {
		t.Fatalf("ClassifierFunc called wrapped function for nil error")
	}
}

func TestClassifierFuncDelegatesNonNilErrors(t *testing.T) {
	errBoom := errors.New("boom")

	var gotErr error
	classifier := ClassifierFunc(func(err error) bool {
		gotErr = err
		return true
	})

	if !classifier.Retryable(errBoom) {
		t.Fatalf("ClassifierFunc.Retryable(non-nil error) = false, want true")
	}
	if !errors.Is(gotErr, errBoom) {
		t.Fatalf("ClassifierFunc received error %v, want %v", gotErr, errBoom)
	}
}

func TestClassifierFuncReturnsWrappedDecision(t *testing.T) {
	errBoom := errors.New("boom")

	tests := []struct {
		name string
		fn   func(error) bool
		want bool
	}{
		{
			name: "wrapped true",
			fn: func(error) bool {
				return true
			},
			want: true,
		},
		{
			name: "wrapped false",
			fn: func(error) bool {
				return false
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := ClassifierFunc(tt.fn)
			if got := classifier.Retryable(errBoom); got != tt.want {
				t.Fatalf("ClassifierFunc.Retryable() = %v, want %v", got, tt.want)
			}
		})
	}
}
