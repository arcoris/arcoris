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

package valuefieldset

// ErrorReason gives stable machine-readable detail within an extraction
// sentinel category.
type ErrorReason string

const (
	// ErrorReasonInvalidZero reports the uninitialized zero value.Value.
	ErrorReasonInvalidZero ErrorReason = "invalid_zero"

	// ErrorReasonInvalidDescriptor reports a malformed or unsupported descriptor.
	ErrorReasonInvalidDescriptor ErrorReason = "invalid_descriptor"

	// ErrorReasonInvalidPath reports a malformed field path.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"

	// ErrorReasonKindMismatch reports a concrete value kind that does not match
	// the descriptor kind.
	ErrorReasonKindMismatch ErrorReason = "kind_mismatch"

	// ErrorReasonUnknownField reports an actual record member rejected by the
	// descriptor's unknown-field policy.
	ErrorReasonUnknownField ErrorReason = "unknown_field"

	// ErrorReasonUnresolvedRef reports a DescriptorRef that cannot be resolved.
	ErrorReasonUnresolvedRef ErrorReason = "unresolved_ref"

	// ErrorReasonReferenceCycle reports recursive DescriptorRef traversal.
	ErrorReasonReferenceCycle ErrorReason = "reference_cycle"

	// ErrorReasonMissingListKey reports an ListMap item missing an
	// identity key field.
	ErrorReasonMissingListKey ErrorReason = "missing_list_key"

	// ErrorReasonInvalidListKey reports an ListMap key value
	// that cannot become a fieldpath selector literal.
	ErrorReasonInvalidListKey ErrorReason = "invalid_list_key"

	// ErrorReasonDuplicateListKey reports repeated ListMap selectors.
	ErrorReasonDuplicateListKey ErrorReason = "duplicate_list_key"

	// ErrorReasonInvalidFieldName reports a record member name that cannot
	// become a fieldpath field element.
	ErrorReasonInvalidFieldName ErrorReason = "invalid_field_name"

	// ErrorReasonInvalidMapKey reports a map member name that cannot become a
	// fieldpath map-key element.
	ErrorReasonInvalidMapKey ErrorReason = "invalid_map_key"
)
