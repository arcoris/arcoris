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

package retrybudget

import (
	"errors"
	"math"
	"testing"
)

func TestNewRatio(t *testing.T) {
	tests := []struct {
		name        string
		numerator   uint64
		denominator uint64
		want        Ratio
		wantErr     error
	}{
		{name: "zero", numerator: 0, denominator: 10, want: RatioZero},
		{name: "one", numerator: 10, denominator: 10, want: RatioOne},
		{name: "reduced", numerator: 2, denominator: 10, want: Ratio{numerator: 1, denominator: 5}},
		{name: "already reduced", numerator: 3, denominator: 7, want: Ratio{numerator: 3, denominator: 7}},
		{name: "zero denominator", numerator: 1, denominator: 0, wantErr: ErrInvalidRatio},
		{name: "above one", numerator: 2, denominator: 1, wantErr: ErrInvalidRatio},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRatio(tt.numerator, tt.denominator)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewRatio() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustRatioPanicsOnInvalidRatio(t *testing.T) {
	requirePanicError(t, ErrInvalidRatio, func() {
		_ = MustRatio(2, 1)
	})
}

func TestRatioAccessors(t *testing.T) {
	ratio := MustRatio(2, 10)
	if got := ratio.Numerator(); got != 1 {
		t.Fatalf("Numerator() = %d, want 1", got)
	}
	if got := ratio.Denominator(); got != 5 {
		t.Fatalf("Denominator() = %d, want 5", got)
	}
}

func TestRatioState(t *testing.T) {
	tests := []struct {
		name      string
		ratio     Ratio
		valid     bool
		zero      bool
		one       bool
		stringVal string
	}{
		{name: "unset", ratio: Ratio{}, valid: false, stringVal: "invalid"},
		{name: "zero", ratio: RatioZero, valid: true, zero: true, stringVal: "0"},
		{name: "one", ratio: RatioOne, valid: true, one: true, stringVal: "1"},
		{name: "fraction", ratio: MustRatio(1, 5), valid: true, stringVal: "1/5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ratio.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
			if got := tt.ratio.IsZero(); got != tt.zero {
				t.Fatalf("IsZero() = %v, want %v", got, tt.zero)
			}
			if got := tt.ratio.IsOne(); got != tt.one {
				t.Fatalf("IsOne() = %v, want %v", got, tt.one)
			}
			if got := tt.ratio.String(); got != tt.stringVal {
				t.Fatalf("String() = %q, want %q", got, tt.stringVal)
			}
		})
	}
}

func TestRatioScaleFloor(t *testing.T) {
	tests := []struct {
		name  string
		ratio Ratio
		value uint64
		want  uint64
	}{
		{name: "zero value", ratio: RatioOne, value: 0, want: 0},
		{name: "zero ratio", ratio: RatioZero, value: 100, want: 0},
		{name: "floor", ratio: MustRatio(1, 5), value: 9, want: 1},
		{name: "exact", ratio: MustRatio(1, 5), value: 10, want: 2},
		{name: "one", ratio: RatioOne, value: 10, want: 10},
		{name: "large exact", ratio: RatioOne, value: math.MaxUint64, want: math.MaxUint64},
		{name: "large half", ratio: MustRatio(1, 2), value: math.MaxUint64, want: math.MaxUint64 / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ratio.ScaleFloor(tt.value); got != tt.want {
				t.Fatalf("ScaleFloor() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRatioScaleFloorPanicsOnInvalidRatio(t *testing.T) {
	requirePanicError(t, ErrInvalidRatio, func() {
		_ = (Ratio{}).ScaleFloor(1)
	})
}

func requirePanicError(t *testing.T, want error, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic = %T(%v), want error %v", recovered, recovered, want)
		}
		if !errors.Is(err, want) {
			t.Fatalf("panic = %v, want %v", err, want)
		}
	}()

	fn()
}
