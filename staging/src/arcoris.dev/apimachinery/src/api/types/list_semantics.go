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

// ListSemantics records future merge/apply intent for list values.
//
// This package only records and validates the descriptor. It does not implement
// patch, apply, field ownership, merge, or conflict behavior.
type ListSemantics uint8

const (
	// ListAtomic treats the complete list as one replaceable value.
	ListAtomic ListSemantics = iota
	// ListOrdered treats physical item indexes as part of the payload contract.
	ListOrdered
	// ListSet treats list elements as set members in future merge layers.
	ListSet
	// ListMap treats object elements as keyed by one or more required fields.
	ListMap
)

// IsValid reports whether s is a known list semantic.
func (s ListSemantics) IsValid() bool {
	return s == ListAtomic || s == ListOrdered || s == ListSet || s == ListMap
}

// String returns a stable diagnostic name for s.
func (s ListSemantics) String() string {
	switch s {
	case ListAtomic:
		return "atomic"
	case ListOrdered:
		return "ordered"
	case ListSet:
		return "set"
	case ListMap:
		return "map"
	default:
		return "unknown"
	}
}
