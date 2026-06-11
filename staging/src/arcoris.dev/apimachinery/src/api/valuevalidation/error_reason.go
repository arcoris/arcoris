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

package valuevalidation

// ErrorReason gives stable machine-readable detail within a validation
// sentinel category.
type ErrorReason string

const (
	// ErrorReasonInvalidZero reports the uninitialized zero value.Value.
	ErrorReasonInvalidZero ErrorReason = "invalid_zero"

	// ErrorReasonInvalidDescriptor reports a malformed or unsupported descriptor.
	ErrorReasonInvalidDescriptor ErrorReason = "invalid_descriptor"

	// ErrorReasonKindMismatch reports a concrete value kind that does not match
	// the descriptor kind.
	ErrorReasonKindMismatch ErrorReason = "kind_mismatch"

	// ErrorReasonNullNotAllowed reports an explicit null where the descriptor is
	// not nullable and is not DescriptorNull.
	ErrorReasonNullNotAllowed ErrorReason = "null_not_allowed"

	// ErrorReasonMissingField reports a required object field that is absent.
	ErrorReasonMissingField ErrorReason = "missing_field"

	// ErrorReasonUnknownField reports an undeclared record member under reject
	// unknown-field policy.
	ErrorReasonUnknownField ErrorReason = "unknown_field"

	// ErrorReasonInvalidFieldName reports a record member name that cannot
	// become a semantic field path element.
	ErrorReasonInvalidFieldName ErrorReason = "invalid_field_name"

	// ErrorReasonInvalidMapKey reports a record member name that cannot become
	// a semantic map-key path element.
	ErrorReasonInvalidMapKey ErrorReason = "invalid_map_key"

	// ErrorReasonBelowMinimum reports a value or length below an inclusive lower
	// bound.
	ErrorReasonBelowMinimum ErrorReason = "below_minimum"

	// ErrorReasonAboveMaximum reports a value or length above an inclusive upper
	// bound.
	ErrorReasonAboveMaximum ErrorReason = "above_maximum"

	// ErrorReasonTooShort reports a string, byte sequence, map, or list shorter
	// than the descriptor minimum.
	ErrorReasonTooShort ErrorReason = "too_short"

	// ErrorReasonTooLong reports a string, byte sequence, map, or list longer
	// than the descriptor maximum.
	ErrorReasonTooLong ErrorReason = "too_long"

	// ErrorReasonPatternMismatch reports a string that does not match its
	// descriptor regexp.
	ErrorReasonPatternMismatch ErrorReason = "pattern_mismatch"

	// ErrorReasonEnumMismatch reports a scalar value outside its descriptor enum.
	ErrorReasonEnumMismatch ErrorReason = "enum_mismatch"

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

	// ErrorReasonDuplicateListSetElement reports repeated ListSet scalar elements.
	ErrorReasonDuplicateListSetElement ErrorReason = "duplicate_list_set_element"
)
