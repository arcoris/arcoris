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

// OrderingMode configures order-sensitive JSON output.
type OrderingMode uint8

const (
	// OrderingDefault defers to the package default during Resolve.
	OrderingDefault OrderingMode = iota

	// OrderingPreserve preserves semantic order supplied to the codec.
	OrderingPreserve

	// OrderingDeterministic sorts deterministic substructures.
	OrderingDeterministic
)

// isKnownOrderingMode reports whether mode is part of the public enum.
func isKnownOrderingMode(mode OrderingMode) bool {
	switch mode {
	case OrderingDefault, OrderingPreserve, OrderingDeterministic:
		return true
	default:
		return false
	}
}
