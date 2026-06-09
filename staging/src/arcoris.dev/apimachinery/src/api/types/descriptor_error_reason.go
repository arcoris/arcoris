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

// DescriptorErrorReason identifies the precise descriptor invariant that failed.
//
// Reason is intentionally separate from Err. Err remains a broad sentinel for
// errors.Is classification, while Reason gives callers stable diagnostic
// detail without parsing human-facing Error text.
type DescriptorErrorReason string

const (
	// DescriptorErrorReasonInvalidDescriptorKind reports an unsupported DescriptorKind value.
	DescriptorErrorReasonInvalidDescriptorKind DescriptorErrorReason = "invalid_descriptor_kind"
	// DescriptorErrorReasonInactivePayload reports a populated payload slot that does not match Descriptor.Code.
	DescriptorErrorReasonInactivePayload DescriptorErrorReason = "inactive_payload"
	// DescriptorErrorReasonInvalidNullability reports a nullability flag that a
	// descriptor kind cannot carry.
	DescriptorErrorReasonInvalidNullability DescriptorErrorReason = "invalid_nullability"

	// DescriptorErrorReasonInvalidRange reports an inverted minimum/maximum rule.
	DescriptorErrorReasonInvalidRange DescriptorErrorReason = "invalid_range"
	// DescriptorErrorReasonNegativeLimit reports a length or size limit below zero.
	DescriptorErrorReasonNegativeLimit DescriptorErrorReason = "negative_limit"
	// DescriptorErrorReasonNonFiniteValue reports NaN or infinity in a float rule.
	DescriptorErrorReasonNonFiniteValue DescriptorErrorReason = "non_finite_value"

	// DescriptorErrorReasonDuplicateEnum reports repeated enum values.
	DescriptorErrorReasonDuplicateEnum DescriptorErrorReason = "duplicate_enum"
	// DescriptorErrorReasonEnumBelowMinimum reports an enum value below a configured minimum.
	DescriptorErrorReasonEnumBelowMinimum DescriptorErrorReason = "enum_below_minimum"
	// DescriptorErrorReasonEnumAboveMaximum reports an enum value above a configured maximum.
	DescriptorErrorReasonEnumAboveMaximum DescriptorErrorReason = "enum_above_maximum"
	// DescriptorErrorReasonEnumPatternMismatch reports an enum value that does not match a string pattern.
	DescriptorErrorReasonEnumPatternMismatch DescriptorErrorReason = "enum_pattern_mismatch"

	// DescriptorErrorReasonInvalidPattern reports a string pattern that cannot be compiled.
	DescriptorErrorReasonInvalidPattern DescriptorErrorReason = "invalid_pattern"
	// DescriptorErrorReasonInvalidPrecision reports an invalid decimal precision rule.
	DescriptorErrorReasonInvalidPrecision DescriptorErrorReason = "invalid_precision"
	// DescriptorErrorReasonInvalidScale reports an invalid decimal scale rule.
	DescriptorErrorReasonInvalidScale DescriptorErrorReason = "invalid_scale"

	// DescriptorErrorReasonMissingElement reports a list descriptor without an element descriptor.
	DescriptorErrorReasonMissingElement DescriptorErrorReason = "missing_element"
	// DescriptorErrorReasonMissingValue reports a map descriptor without a value descriptor.
	DescriptorErrorReasonMissingValue DescriptorErrorReason = "missing_value"
	// DescriptorErrorReasonInvalidSemantics reports an unsupported list semantic policy.
	DescriptorErrorReasonInvalidSemantics DescriptorErrorReason = "invalid_semantics"

	// DescriptorErrorReasonInvalidFieldName reports a malformed object field name.
	DescriptorErrorReasonInvalidFieldName DescriptorErrorReason = "invalid_field_name"
	// DescriptorErrorReasonDuplicateFieldName reports a repeated field name in one object descriptor.
	DescriptorErrorReasonDuplicateFieldName DescriptorErrorReason = "duplicate_field_name"
	// DescriptorErrorReasonInvalidPresence reports a field without Required or Optional presence.
	DescriptorErrorReasonInvalidPresence DescriptorErrorReason = "invalid_presence"
	// DescriptorErrorReasonInvalidUnknownPolicy reports an unsupported object unknown-field policy.
	DescriptorErrorReasonInvalidUnknownPolicy DescriptorErrorReason = "invalid_unknown_policy"

	// DescriptorErrorReasonInvalidReferenceName reports a malformed DescriptorRef name.
	DescriptorErrorReasonInvalidReferenceName DescriptorErrorReason = "invalid_reference_name"
	// DescriptorErrorReasonUnknownReference reports a DescriptorRef that a Resolver cannot resolve.
	DescriptorErrorReasonUnknownReference DescriptorErrorReason = "unknown_reference"
	// DescriptorErrorReasonReferenceCycle reports a recursive Definition graph.
	DescriptorErrorReasonReferenceCycle DescriptorErrorReason = "reference_cycle"
	// DescriptorErrorReasonInvalidResolvedDefinition reports a DescriptorRef target that
	// resolves but is structurally invalid.
	DescriptorErrorReasonInvalidResolvedDefinition DescriptorErrorReason = "invalid_resolved_definition"

	// DescriptorErrorReasonMissingListMapKey reports ListMap semantics without map keys.
	DescriptorErrorReasonMissingListMapKey DescriptorErrorReason = "missing_list_map_key"
	// DescriptorErrorReasonInvalidListMapKey reports a malformed ListMap key name.
	DescriptorErrorReasonInvalidListMapKey DescriptorErrorReason = "invalid_list_map_key"
	// DescriptorErrorReasonDuplicateListMapKey reports a repeated ListMap key name.
	DescriptorErrorReasonDuplicateListMapKey DescriptorErrorReason = "duplicate_list_map_key"
	// DescriptorErrorReasonListMapKeyNotFound reports a ListMap key absent from the object element.
	DescriptorErrorReasonListMapKeyNotFound DescriptorErrorReason = "list_map_key_not_found"
	// DescriptorErrorReasonListMapKeyOptional reports a ListMap key field that is not required.
	DescriptorErrorReasonListMapKeyOptional DescriptorErrorReason = "list_map_key_optional"
	// DescriptorErrorReasonInvalidListMapKeyDescriptor reports a ListMap key field that
	// cannot produce stable selector identity.
	DescriptorErrorReasonInvalidListMapKeyDescriptor DescriptorErrorReason = "invalid_list_map_key_descriptor"
	// DescriptorErrorReasonListMapElementNotObject reports ListMap semantics over a non-object element.
	DescriptorErrorReasonListMapElementNotObject DescriptorErrorReason = "list_map_element_not_object"
	// DescriptorErrorReasonInvalidListSetElement reports a ListSet element descriptor
	// that cannot produce stable item identity.
	DescriptorErrorReasonInvalidListSetElement DescriptorErrorReason = "invalid_list_set_element"

	// DescriptorErrorReasonInvalidMapKey reports an unsupported dynamic-map key descriptor.
	DescriptorErrorReasonInvalidMapKey DescriptorErrorReason = "invalid_map_key"
)
