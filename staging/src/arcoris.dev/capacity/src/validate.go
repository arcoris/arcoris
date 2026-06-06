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

package capacity

// requirePositiveAmount panics when amount cannot represent a real reservation.
func requirePositiveAmount(amount Amount) {
	if amount.IsZero() {
		panicAt("amount", ErrZeroAmount, "amount must be positive")
	}
}

// requireValidVector panics when vector is not canonical.
func requireValidVector(path string, vector Vector) {
	if !vector.IsValid() {
		panicAt(path, ErrInvalidVector, "vector must be canonical")
	}
}

// requireValidDemand panics when demand is empty or not canonical.
func requireValidDemand(path string, demand Demand) {
	if !demand.IsValid() {
		panicAt(path, ErrInvalidDemand, "demand must be non-empty and canonical")
	}
}
