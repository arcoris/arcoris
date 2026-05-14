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

package result_test

import (
	"errors"
	"testing"

	"arcoris.dev/value/result"
)

func TestZeroValueIsOK(t *testing.T) {
	var r result.Result[string]

	if !r.IsOK() {
		t.Fatal("zero value must be OK")
	}
	if r.IsErr() {
		t.Fatal("zero value must not be Err")
	}

	got, err := r.Load()
	if err != nil {
		t.Fatalf("Load returned error for zero value: %v", err)
	}
	if got != "" {
		t.Fatalf("Load returned non-zero value for zero Result: %q", got)
	}
}

func TestOK(t *testing.T) {
	r := result.OK("value")

	if !r.IsOK() {
		t.Fatal("OK value must report IsOK")
	}
	if r.IsErr() {
		t.Fatal("OK value must not report IsErr")
	}

	got, err := r.Load()
	if err != nil {
		t.Fatalf("Load returned error for OK: %v", err)
	}
	if got != "value" {
		t.Fatalf("got %q, want %q", got, "value")
	}
}

func TestErr(t *testing.T) {
	want := errors.New("failed")
	r := result.Err[string](want)

	if !r.IsErr() {
		t.Fatal("Err value must report IsErr")
	}
	if r.IsOK() {
		t.Fatal("Err value must not report IsOK")
	}
	if got := r.Err(); !errors.Is(got, want) {
		t.Fatalf("Err() = %v, want %v", got, want)
	}

	got, err := r.Load()
	if !errors.Is(err, want) {
		t.Fatalf("Load error = %v, want %v", err, want)
	}
	if got != "" {
		t.Fatalf("Load returned non-zero value for Err: %q", got)
	}
}

func TestErrPanicsOnNilError(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Err did not panic for nil error")
		}
	}()

	_ = result.Err[string](nil)
}

func TestFrom(t *testing.T) {
	boom := errors.New("boom")
	tests := []struct {
		name  string
		err   error
		isErr bool
	}{
		{name: "ok", err: nil, isErr: false},
		{name: "err", err: boom, isErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := result.From("value", tc.err)
			if got := r.IsErr(); got != tc.isErr {
				t.Fatalf("IsErr() = %v, want %v", got, tc.isErr)
			}
		})
	}
}

func TestValueOr(t *testing.T) {
	if got := result.OK("value").ValueOr("fallback"); got != "value" {
		t.Fatalf("OK.ValueOr() = %q, want %q", got, "value")
	}
	if got := result.Err[string](errors.New("failed")).ValueOr("fallback"); got != "fallback" {
		t.Fatalf("Err.ValueOr() = %q, want %q", got, "fallback")
	}
}

func TestMust(t *testing.T) {
	if got := result.OK("value").Must(); got != "value" {
		t.Fatalf("Must() = %q, want %q", got, "value")
	}
}

func TestMustPanicsForErr(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Must did not panic for Err")
		}
	}()

	_ = result.Err[string](errors.New("failed")).Must()
}
