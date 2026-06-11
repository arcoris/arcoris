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

package valuemerge

// ErrorReason gives stable machine-readable detail inside a broad category.
type ErrorReason string

const (
	// ErrorReasonInvalidZero reports the uninitialized zero value.Value.
	ErrorReasonInvalidZero ErrorReason = "invalid_zero"

	// ErrorReasonInvalidDescriptor reports a malformed or unsupported descriptor.
	ErrorReasonInvalidDescriptor ErrorReason = "invalid_descriptor"

	// ErrorReasonInvalidPath reports a malformed semantic field path.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"

	// ErrorReasonKindMismatch reports descriptor/value kind incompatibility.
	ErrorReasonKindMismatch ErrorReason = "kind_mismatch"

	// ErrorReasonUnknownField reports an undeclared record member rejected by policy.
	ErrorReasonUnknownField ErrorReason = "unknown_field"

	// ErrorReasonInvalidFieldName reports a payload name that cannot become a field element.
	ErrorReasonInvalidFieldName ErrorReason = "invalid_field_name"

	// ErrorReasonInvalidMapKey reports a payload name that cannot become a map-key element.
	ErrorReasonInvalidMapKey ErrorReason = "invalid_map_key"

	// ErrorReasonUnresolvedRef reports a DescriptorRef the resolver cannot load.
	ErrorReasonUnresolvedRef ErrorReason = "unresolved_ref"

	// ErrorReasonReferenceCycle reports recursive or over-depth DescriptorRef traversal.
	ErrorReasonReferenceCycle ErrorReason = "reference_cycle"

	// ErrorReasonMissingListKey reports a ListMap item missing a key field.
	ErrorReasonMissingListKey ErrorReason = "missing_list_key"

	// ErrorReasonInvalidListKey reports a ListMap key that cannot form a selector.
	ErrorReasonInvalidListKey ErrorReason = "invalid_list_key"

	// ErrorReasonDuplicateListKey reports repeated ListMap selector identity.
	ErrorReasonDuplicateListKey ErrorReason = "duplicate_list_key"

	// ErrorReasonUnsupportedMerge reports a selected path blocked by semantics.
	ErrorReasonUnsupportedMerge ErrorReason = "unsupported_merge"
)
