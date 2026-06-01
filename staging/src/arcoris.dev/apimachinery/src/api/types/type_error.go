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

import (
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidType classifies malformed structural type descriptors.
	ErrInvalidType = errors.New("invalid type")
	// ErrInvalidTypeCode classifies unsupported TypeCode values.
	ErrInvalidTypeCode = errors.New("invalid type code")
	// ErrInvalidField classifies malformed object field descriptors.
	ErrInvalidField = errors.New("invalid field")
	// ErrDuplicateField classifies repeated field names within one object.
	ErrDuplicateField = errors.New("duplicate field")
	// ErrInvalidTypeReference classifies malformed TypeRef names or ref cycles.
	ErrInvalidTypeReference = errors.New("invalid type reference")
	// ErrUnknownTypeReference classifies TypeRef names absent from a Resolver.
	ErrUnknownTypeReference = errors.New("unknown type reference")
)

// TypeErrorReason identifies the precise descriptor invariant that failed.
//
// Reason is intentionally separate from Err. Err remains a broad sentinel for
// errors.Is classification, while Reason gives callers stable diagnostic
// detail without parsing human-facing Error text.
type TypeErrorReason string

const (
	// TypeErrorReasonInvalidTypeCode reports an unsupported TypeCode value.
	TypeErrorReasonInvalidTypeCode TypeErrorReason = "invalid_type_code"
	// TypeErrorReasonInactivePayload reports a populated payload slot that does not match Type.Code.
	TypeErrorReasonInactivePayload TypeErrorReason = "inactive_payload"
	// TypeErrorReasonInvalidNullability reports a nullability flag that a descriptor kind cannot carry.
	TypeErrorReasonInvalidNullability TypeErrorReason = "invalid_nullability"

	// TypeErrorReasonInvalidRange reports an inverted minimum/maximum rule.
	TypeErrorReasonInvalidRange TypeErrorReason = "invalid_range"
	// TypeErrorReasonNegativeLimit reports a length or size limit below zero.
	TypeErrorReasonNegativeLimit TypeErrorReason = "negative_limit"
	// TypeErrorReasonNonFiniteValue reports NaN or infinity in a float rule.
	TypeErrorReasonNonFiniteValue TypeErrorReason = "non_finite_value"

	// TypeErrorReasonDuplicateEnum reports repeated enum values.
	TypeErrorReasonDuplicateEnum TypeErrorReason = "duplicate_enum"
	// TypeErrorReasonEnumBelowMinimum reports an enum value below a configured minimum.
	TypeErrorReasonEnumBelowMinimum TypeErrorReason = "enum_below_minimum"
	// TypeErrorReasonEnumAboveMaximum reports an enum value above a configured maximum.
	TypeErrorReasonEnumAboveMaximum TypeErrorReason = "enum_above_maximum"
	// TypeErrorReasonEnumPatternMismatch reports an enum value that does not match a string pattern.
	TypeErrorReasonEnumPatternMismatch TypeErrorReason = "enum_pattern_mismatch"

	// TypeErrorReasonInvalidPattern reports a string pattern that cannot be compiled.
	TypeErrorReasonInvalidPattern TypeErrorReason = "invalid_pattern"
	// TypeErrorReasonInvalidPrecision reports an invalid decimal precision rule.
	TypeErrorReasonInvalidPrecision TypeErrorReason = "invalid_precision"
	// TypeErrorReasonInvalidScale reports an invalid decimal scale rule.
	TypeErrorReasonInvalidScale TypeErrorReason = "invalid_scale"

	// TypeErrorReasonMissingElement reports a list descriptor without an element type.
	TypeErrorReasonMissingElement TypeErrorReason = "missing_element"
	// TypeErrorReasonMissingValue reports a map descriptor without a value type.
	TypeErrorReasonMissingValue TypeErrorReason = "missing_value"
	// TypeErrorReasonInvalidSemantics reports an unsupported list semantic policy.
	TypeErrorReasonInvalidSemantics TypeErrorReason = "invalid_semantics"

	// TypeErrorReasonInvalidFieldName reports a malformed object field name.
	TypeErrorReasonInvalidFieldName TypeErrorReason = "invalid_field_name"
	// TypeErrorReasonDuplicateFieldName reports a repeated field name in one object descriptor.
	TypeErrorReasonDuplicateFieldName TypeErrorReason = "duplicate_field_name"
	// TypeErrorReasonInvalidPresence reports a field without Required or Optional presence.
	TypeErrorReasonInvalidPresence TypeErrorReason = "invalid_presence"
	// TypeErrorReasonInvalidUnknownPolicy reports an unsupported object unknown-field policy.
	TypeErrorReasonInvalidUnknownPolicy TypeErrorReason = "invalid_unknown_policy"

	// TypeErrorReasonInvalidReferenceName reports a malformed TypeRef name.
	TypeErrorReasonInvalidReferenceName TypeErrorReason = "invalid_reference_name"
	// TypeErrorReasonUnknownReference reports a TypeRef that a Resolver cannot resolve.
	TypeErrorReasonUnknownReference TypeErrorReason = "unknown_reference"
	// TypeErrorReasonReferenceCycle reports a recursive TypeDefinition graph.
	TypeErrorReasonReferenceCycle TypeErrorReason = "reference_cycle"
	// TypeErrorReasonInvalidResolvedDefinition reports a TypeRef target that resolves but is structurally invalid.
	TypeErrorReasonInvalidResolvedDefinition TypeErrorReason = "invalid_resolved_definition"

	// TypeErrorReasonMissingListMapKey reports ListMap semantics without map keys.
	TypeErrorReasonMissingListMapKey TypeErrorReason = "missing_list_map_key"
	// TypeErrorReasonInvalidListMapKey reports a malformed ListMap key name.
	TypeErrorReasonInvalidListMapKey TypeErrorReason = "invalid_list_map_key"
	// TypeErrorReasonListMapKeyNotFound reports a ListMap key absent from the object element.
	TypeErrorReasonListMapKeyNotFound TypeErrorReason = "list_map_key_not_found"
	// TypeErrorReasonListMapKeyOptional reports a ListMap key field that is not required.
	TypeErrorReasonListMapKeyOptional TypeErrorReason = "list_map_key_optional"
	// TypeErrorReasonInvalidListMapKeyType reports a ListMap key field that cannot produce stable selector identity.
	TypeErrorReasonInvalidListMapKeyType TypeErrorReason = "invalid_list_map_key_type"
	// TypeErrorReasonListMapElementNotObject reports ListMap semantics over a non-object element.
	TypeErrorReasonListMapElementNotObject TypeErrorReason = "list_map_element_not_object"

	// TypeErrorReasonInvalidMapKey reports an unsupported dynamic-map key type.
	TypeErrorReasonInvalidMapKey TypeErrorReason = "invalid_map_key"
)

// TypeError attaches structured descriptor diagnostics to a classified error.
//
// Path is a descriptor path such as object.fields[spec].type, list.elem, or
// ref(arcoris.meta.Name). It is not a path into a future concrete API object.
type TypeError struct {
	// Record stores the shared path, sentinel, reason, and detail fields.
	diagnostic.Record[TypeErrorReason]
}

// Error returns a stable diagnostic message for e.
func (e *TypeError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Record.Format("types")
}

// Unwrap returns the classified validation error.
func (e *TypeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Record.Unwrap()
}

// typeError creates a path-aware validation error.
func typeError(path string, err error) error {
	return &TypeError{
		Record: diagnostic.NewRecord(path, err, TypeErrorReason(""), ""),
	}
}

// typeErrorf creates a path-aware validation error with structured detail.
func typeErrorf(path string, err error, reason TypeErrorReason, format string, args ...any) error {
	detail := ""
	if format != "" {
		detail = fmt.Sprintf(format, args...)
	}
	return &TypeError{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}
