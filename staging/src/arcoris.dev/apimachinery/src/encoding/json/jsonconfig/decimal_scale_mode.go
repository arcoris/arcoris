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

// DecimalScaleMode controls decimal scale rendering.
type DecimalScaleMode uint8

const (
	// DecimalScaleDefault defers to the package default during Resolve.
	DecimalScaleDefault DecimalScaleMode = iota

	// DecimalScalePreserve writes decimal text exactly as stored by value.Decimal.
	DecimalScalePreserve

	// DecimalScaleTrimTrailingZeros is reserved for canonical decimal profiles.
	DecimalScaleTrimTrailingZeros
)

// isKnownDecimalScaleMode reports whether mode is part of the public enum.
func isKnownDecimalScaleMode(mode DecimalScaleMode) bool {
	switch mode {
	case DecimalScaleDefault, DecimalScalePreserve, DecimalScaleTrimTrailingZeros:
		return true
	default:
		return false
	}
}
