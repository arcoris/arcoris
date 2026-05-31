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

package value

import (
	"math"
	"testing"
)

func TestNewIntegerConstructors(t *testing.T) {
	tests := []struct {
		name      string
		integer   Integer
		negative  bool
		magnitude uint64
	}{
		{name: "positive int64", integer: NewIntegerFromInt64(42), magnitude: 42},
		{name: "negative int64", integer: NewIntegerFromInt64(-42), negative: true, magnitude: 42},
		{
			name:      "min int64",
			integer:   NewIntegerFromInt64(math.MinInt64),
			negative:  true,
			magnitude: uint64(math.MaxInt64) + 1,
		},
		{name: "max uint64", integer: NewIntegerFromUint64(math.MaxUint64), magnitude: math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.integer.IsNegative(), tt.negative)
			requireEqual(t, tt.integer.Magnitude(), tt.magnitude)
		})
	}
}
