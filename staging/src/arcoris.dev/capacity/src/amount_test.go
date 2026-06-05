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

package capacity_test

import (
	"math"
	"testing"

	"arcoris.dev/capacity"
)

func TestAmountPredicatesAndConversion(t *testing.T) {
	t.Parallel()

	if !capacity.Amount(0).IsZero() {
		t.Fatal("zero amount was not zero")
	}
	if capacity.Amount(0).IsPositive() {
		t.Fatal("zero amount was positive")
	}
	if !capacity.Amount(1).IsPositive() {
		t.Fatal("positive amount was not positive")
	}
	if got := capacity.Amount(42).Uint64(); got != 42 {
		t.Fatalf("Uint64() = %d, want 42", got)
	}
}

func TestAmountCompare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    capacity.Amount
		b    capacity.Amount
		want int
	}{
		{name: "less", a: 1, b: 2, want: -1},
		{name: "equal", a: 2, b: 2, want: 0},
		{name: "greater", a: 3, b: 2, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.a.Compare(tt.b); got != tt.want {
				t.Fatalf("Compare() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAmountCheckedArithmetic(t *testing.T) {
	t.Parallel()

	if got, ok := capacity.Amount(2).CheckedAdd(3); !ok || got != 5 {
		t.Fatalf("CheckedAdd() = %d, %v; want 5, true", got, ok)
	}
	if _, ok := capacity.Amount(math.MaxUint64).CheckedAdd(1); ok {
		t.Fatal("CheckedAdd() overflow returned ok=true")
	}
	if got, ok := capacity.Amount(5).CheckedSub(3); !ok || got != 2 {
		t.Fatalf("CheckedSub() = %d, %v; want 2, true", got, ok)
	}
	if _, ok := capacity.Amount(1).CheckedSub(2); ok {
		t.Fatal("CheckedSub() underflow returned ok=true")
	}
	if got := capacity.Amount(1).SaturatingSub(2); got != 0 {
		t.Fatalf("SaturatingSub() = %d, want 0", got)
	}
}
