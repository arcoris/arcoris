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

package types

// UnknownFieldPolicy records how future value processors should handle fields
// not declared by an object descriptor.
//
// The policy is structural metadata only. This package validates the policy but
// does not reject, prune, preserve, or otherwise process concrete object values.
type UnknownFieldPolicy uint8

const (
	// UnknownReject means undeclared object fields are rejected by value layers.
	UnknownReject UnknownFieldPolicy = iota
	// UnknownPrune means undeclared object fields may be dropped by value layers.
	UnknownPrune
	// UnknownPreserveOpaque means undeclared object fields may be retained by
	// value layers as opaque payload.
	//
	// The descriptor does not provide structural traversal for those fields.
	UnknownPreserveOpaque
)

// IsValid reports whether p is a known unknown-field policy.
func (p UnknownFieldPolicy) IsValid() bool {
	return p == UnknownReject || p == UnknownPrune || p == UnknownPreserveOpaque
}
