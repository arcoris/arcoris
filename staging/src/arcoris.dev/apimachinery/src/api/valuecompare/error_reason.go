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

package valuecompare

// ErrorReason gives stable machine-readable detail within a comparison
// sentinel category.
type ErrorReason string

const (
	// ErrorReasonInvalidZero identifies the invalid zero value.Value sentinel.
	ErrorReasonInvalidZero ErrorReason = "invalid_zero"

	// ErrorReasonInvalidDescriptor identifies unusable descriptor views or codes.
	ErrorReasonInvalidDescriptor ErrorReason = "invalid_descriptor"

	// ErrorReasonKindMismatch identifies descriptor/value kind incompatibility.
	ErrorReasonKindMismatch ErrorReason = "kind_mismatch"

	// ErrorReasonUnknownField identifies a rejected object member.
	ErrorReasonUnknownField ErrorReason = "unknown_field"

	// ErrorReasonUnresolvedRef identifies a TypeRef the resolver cannot load.
	ErrorReasonUnresolvedRef ErrorReason = "unresolved_ref"

	// ErrorReasonReferenceCycle identifies recursive TypeRef traversal.
	ErrorReasonReferenceCycle ErrorReason = "reference_cycle"

	// ErrorReasonMissingListKey identifies a ListMap item missing a key field.
	ErrorReasonMissingListKey ErrorReason = "missing_list_key"

	// ErrorReasonInvalidListKey identifies a ListMap key that cannot form a selector.
	ErrorReasonInvalidListKey ErrorReason = "invalid_list_key"

	// ErrorReasonDuplicateListKey identifies repeated ListMap selector identity.
	ErrorReasonDuplicateListKey ErrorReason = "duplicate_list_key"
)
