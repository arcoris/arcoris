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

package valuevalidation

// signedWidthLimits returns the fixed-width bounds for a signed descriptor.
func signedWidthLimits(lowerValue, upperValue int64) integerLimits[int64] {
	return integerLimits[int64]{
		lower: integerBound[int64]{value: lowerValue, set: true},
		upper: integerBound[int64]{value: upperValue, set: true},
	}
}

// unsignedWidthLimits returns the fixed-width upper bound for an unsigned descriptor.
func unsignedWidthLimits(upperValue uint64) integerLimits[uint64] {
	return integerLimits[uint64]{
		upper: integerBound[uint64]{value: upperValue, set: true},
	}
}

// exactIntegerLimits reads bounds from a descriptor view that already uses the
// normalized validation type.
func exactIntegerLimits[T integerValue](
	lowerAccessor func() (T, bool),
	upperAccessor func() (T, bool),
) integerLimits[T] {
	lowerValue, lowerSet := lowerAccessor()
	upperValue, upperSet := upperAccessor()

	return integerLimits[T]{
		lower: integerBound[T]{value: lowerValue, set: lowerSet},
		upper: integerBound[T]{value: upperValue, set: upperSet},
	}
}

// signedDescriptorLimits widens small signed descriptor limits into int64.
func signedDescriptorLimits[T ~int8 | ~int16 | ~int32](
	lowerAccessor func() (T, bool),
	upperAccessor func() (T, bool),
) integerLimits[int64] {
	lowerValue, lowerSet := lowerAccessor()
	upperValue, upperSet := upperAccessor()

	return integerLimits[int64]{
		lower: integerBound[int64]{value: int64(lowerValue), set: lowerSet},
		upper: integerBound[int64]{value: int64(upperValue), set: upperSet},
	}
}

// unsignedDescriptorLimits widens small unsigned descriptor limits into uint64.
func unsignedDescriptorLimits[T ~uint8 | ~uint16 | ~uint32](
	lowerAccessor func() (T, bool),
	upperAccessor func() (T, bool),
) integerLimits[uint64] {
	lowerValue, lowerSet := lowerAccessor()
	upperValue, upperSet := upperAccessor()

	return integerLimits[uint64]{
		lower: integerBound[uint64]{value: uint64(lowerValue), set: lowerSet},
		upper: integerBound[uint64]{value: uint64(upperValue), set: upperSet},
	}
}
