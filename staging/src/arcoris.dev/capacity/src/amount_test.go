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
	"testing"

	"arcoris.dev/capacity"
)

func TestAmountPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   capacity.Amount
		zero     bool
		positive bool
		raw      uint64
	}{
		{name: "zero", amount: 0, zero: true, positive: false, raw: 0},
		{name: "one", amount: 1, zero: false, positive: true, raw: 1},
		{name: "large", amount: capacity.Amount(^uint64(0)), zero: false, positive: true, raw: ^uint64(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.amount.IsZero(); got != tt.zero {
				t.Fatalf("IsZero() = %t, want %t", got, tt.zero)
			}
			if got := tt.amount.IsPositive(); got != tt.positive {
				t.Fatalf("IsPositive() = %t, want %t", got, tt.positive)
			}
			if got := tt.amount.Uint64(); got != tt.raw {
				t.Fatalf("Uint64() = %d, want %d", got, tt.raw)
			}
		})
	}
}
