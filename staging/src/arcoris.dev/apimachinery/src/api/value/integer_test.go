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

func TestIntegerFitsAndAccessors(t *testing.T) {
	min := NewIntegerFromInt64(math.MinInt64)
	gotInt64, ok := min.Int64()
	requireEqual(t, ok, true)
	requireEqual(t, gotInt64, int64(math.MinInt64))
	requireEqual(t, min.FitsUint64(), false)

	maxUint := NewIntegerFromUint64(math.MaxUint64)
	_, ok = maxUint.Int64()
	requireEqual(t, ok, false)
	requireEqual(t, maxUint.FitsInt64(), false)

	gotUint64, ok := maxUint.Uint64()
	requireEqual(t, ok, true)
	requireEqual(t, gotUint64, uint64(math.MaxUint64))
}

func TestIntegerString(t *testing.T) {
	requireEqual(t, NewIntegerFromInt64(-42).String(), "-42")
	requireEqual(t, NewIntegerFromUint64(math.MaxUint64).String(), "18446744073709551615")
}

func TestIntegerCompareAndEqual(t *testing.T) {
	tests := []struct {
		name string
		a    Integer
		b    Integer
		want int
	}{
		{name: "less negative", a: NewIntegerFromInt64(-2), b: NewIntegerFromInt64(-1), want: -1},
		{name: "negative less than positive", a: NewIntegerFromInt64(-1), b: NewIntegerFromUint64(1), want: -1},
		{name: "equal", a: NewIntegerFromInt64(7), b: NewIntegerFromUint64(7), want: 0},
		{name: "greater", a: NewIntegerFromUint64(9), b: NewIntegerFromInt64(7), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.a.Compare(tt.b), tt.want)
			requireEqual(t, tt.a.Equal(tt.b), tt.want == 0)
		})
	}
}
