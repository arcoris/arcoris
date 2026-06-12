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

// OwnershipNormalizeMode controls ownership state normalization before encoding.
type OwnershipNormalizeMode uint8

const (
	// OwnershipNormalizeDefault defers to the package default during Resolve.
	OwnershipNormalizeDefault OwnershipNormalizeMode = iota

	// OwnershipNormalizeNever skips pre-encode normalization.
	OwnershipNormalizeNever

	// OwnershipNormalizeWhenDeterministic normalizes only with deterministic output.
	OwnershipNormalizeWhenDeterministic

	// OwnershipNormalizeAlways always normalizes before encoding.
	OwnershipNormalizeAlways
)

// isKnownOwnershipNormalizeMode reports whether mode is part of the public enum.
func isKnownOwnershipNormalizeMode(mode OwnershipNormalizeMode) bool {
	switch mode {
	case OwnershipNormalizeDefault,
		OwnershipNormalizeNever,
		OwnershipNormalizeWhenDeterministic,
		OwnershipNormalizeAlways:
		return true
	default:
		return false
	}
}
