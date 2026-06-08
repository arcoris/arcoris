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

package fixedwindow

import "arcoris.dev/resilience/retrybudget"

// Ratio is the exact fixed-window retry allowance ratio.
//
// It aliases retrybudget.Ratio so snapshots, configuration, and tests all share
// one domain value instead of converting through floating point.
type Ratio = retrybudget.Ratio

var (
	// RatioZero is the exact 0% retry allowance ratio.
	RatioZero = retrybudget.RatioZero

	// RatioOne is the exact 100% retry allowance ratio.
	RatioOne = retrybudget.RatioOne
)

// NewRatio returns the canonical exact ratio numerator/denominator.
func NewRatio(numerator, denominator uint64) (Ratio, error) {
	return retrybudget.NewRatio(numerator, denominator)
}

// MustRatio returns the canonical exact ratio numerator/denominator.
//
// MustRatio panics with retrybudget.ErrInvalidRatio when NewRatio would return
// an error.
func MustRatio(numerator, denominator uint64) Ratio {
	return retrybudget.MustRatio(numerator, denominator)
}
