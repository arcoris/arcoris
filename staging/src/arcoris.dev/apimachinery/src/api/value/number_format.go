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

import "strconv"

// appendPaddedUnsignedDecimal appends value with zero padding up to width.
//
// The helper is shared by temporal diagnostics, where the value model wants
// simple fixed-width text without importing fmt.
func appendPaddedUnsignedDecimal(out []byte, value uint64, width int) []byte {
	digits := strconv.FormatUint(value, 10)
	for i := len(digits); i < width; i++ {
		out = append(out, '0')
	}

	return append(out, digits...)
}

// appendPaddedSignedDecimal appends value with fmt-like zero padding.
//
// For negative values, width includes the minus sign, matching %0*d behavior.
func appendPaddedSignedDecimal(out []byte, value int, width int) []byte {
	if value >= 0 {
		return appendPaddedUnsignedDecimal(out, uint64(value), width)
	}

	out = append(out, '-')
	return appendPaddedUnsignedDecimal(out, signedMagnitude(value), width-1)
}

// signedMagnitude returns abs(value) without overflowing on the minimum int.
func signedMagnitude(value int) uint64 {
	n := int64(value)
	return uint64(-(n + 1)) + 1
}
