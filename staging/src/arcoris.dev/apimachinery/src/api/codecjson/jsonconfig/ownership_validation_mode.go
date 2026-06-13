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

// OwnershipValidationMode controls objectownership.State validation after decode.
type OwnershipValidationMode uint8

const (
	// OwnershipValidationDefault defers to the package default during Resolve.
	OwnershipValidationDefault OwnershipValidationMode = iota

	// OwnershipValidationEnable validates decoded ownership state.
	OwnershipValidationEnable

	// OwnershipValidationDisable skips decoded ownership state validation.
	OwnershipValidationDisable
)

// isKnownOwnershipValidationMode reports whether mode is part of the public enum.
func isKnownOwnershipValidationMode(mode OwnershipValidationMode) bool {
	switch mode {
	case OwnershipValidationDefault, OwnershipValidationEnable, OwnershipValidationDisable:
		return true
	default:
		return false
	}
}
