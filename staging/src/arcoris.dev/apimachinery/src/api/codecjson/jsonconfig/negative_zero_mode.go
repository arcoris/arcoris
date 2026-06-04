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

package jsonconfig

// NegativeZeroMode controls negative floating zero encoding.
type NegativeZeroMode uint8

const (
	// NegativeZeroDefault defers to the package default during Resolve.
	NegativeZeroDefault NegativeZeroMode = iota

	// NegativeZeroNormalize writes negative floating zero as 0.
	NegativeZeroNormalize

	// NegativeZeroPreserve writes negative floating zero as -0.
	NegativeZeroPreserve

	// NegativeZeroReject rejects negative floating zero.
	NegativeZeroReject
)

// isKnownNegativeZeroMode reports whether mode is part of the public enum.
func isKnownNegativeZeroMode(mode NegativeZeroMode) bool {
	switch mode {
	case NegativeZeroDefault, NegativeZeroNormalize, NegativeZeroPreserve, NegativeZeroReject:
		return true
	default:
		return false
	}
}
