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

// NewIntegerFromInt64 constructs an Integer from a signed 64-bit value.
//
// math.MinInt64 is handled without negating it directly, because -MinInt64
// overflows int64. The resulting Integer remains exact.
func NewIntegerFromInt64(v int64) Integer {
	if v >= 0 {
		return Integer{magnitude: uint64(v)}
	}

	magnitude := uint64(-(v + 1)) + 1
	return Integer{negative: true, magnitude: magnitude}
}

// NewIntegerFromUint64 constructs an Integer from an unsigned 64-bit value.
//
// The unsigned domain is accepted in full, including math.MaxUint64. Descriptor
// validation can later decide whether a smaller signed or unsigned width is
// allowed for a particular field.
func NewIntegerFromUint64(v uint64) Integer {
	return Integer{magnitude: v}
}
