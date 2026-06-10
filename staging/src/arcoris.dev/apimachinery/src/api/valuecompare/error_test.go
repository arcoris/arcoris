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

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"errors"
	"testing"
)

func TestErrorIsSentinel(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "bad")

	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("errors.Is(err, ErrInvalidValue) = false")
	}
}

func TestErrorIsInvalidPath(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidPath, ErrorReasonInvalidPath, "bad")

	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("errors.Is(err, ErrInvalidPath) = false")
	}
}

func TestErrorAsValueCompareError(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "bad")

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if got.Path != "$" {
		t.Fatalf("path = %q, want $", got.Path)
	}
}

func TestNilErrorStringAndUnwrap(t *testing.T) {
	var err *Error

	if err.Error() != "<nil>" {
		t.Fatalf("nil Error() = %q", err.Error())
	}
	if err.Unwrap() != nil {
		t.Fatalf("nil Unwrap() != nil")
	}
}
func TestErrorAtBuildsStructuredError(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "bad")

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidZero)
	requireErrorPath(t, err, "$")
	requireErrorDetailContains(t, err, "bad")
}

func TestErrorfAtBuildsFormattedDetail(t *testing.T) {
	err := errorfAt(fieldpath.Root(), ErrUnknownField, ErrorReasonUnknownField, "field %q", "extra")

	requireErrorDetailContains(t, err, `field "extra"`)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt(fieldpath.Root(), ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "bad", cause)

	requireErrorIs(t, err, ErrInvalidDescriptor)
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(err, cause) = false")
	}
}
func TestErrorReasonStrings(t *testing.T) {
	cases := map[ErrorReason]string{
		ErrorReasonInvalidZero:       "invalid_zero",
		ErrorReasonInvalidDescriptor: "invalid_descriptor",
		ErrorReasonInvalidPath:       "invalid_path",
		ErrorReasonKindMismatch:      "kind_mismatch",
		ErrorReasonUnknownField:      "unknown_field",
		ErrorReasonUnresolvedRef:     "unresolved_ref",
		ErrorReasonReferenceCycle:    "reference_cycle",
		ErrorReasonMissingListKey:    "missing_list_key",
		ErrorReasonInvalidListKey:    "invalid_list_key",
		ErrorReasonDuplicateListKey:  "duplicate_list_key",
	}

	for reason, want := range cases {
		if string(reason) != want {
			t.Fatalf("reason = %q, want %q", reason, want)
		}
	}
}

func TestErrorReasonInvalidPath(t *testing.T) {
	if string(ErrorReasonInvalidPath) != "invalid_path" {
		t.Fatalf("invalid path reason = %q", ErrorReasonInvalidPath)
	}
}
func TestAddedSubtreeErrorUsesValueCompareErrorModel(t *testing.T) {
	descriptor := types.Object(types.Field("name").String().Optional()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("name", value.BoolValue(true)))

	_, err := Compare(valueObject(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrKindMismatch)
	requireErrorReason(t, err, ErrorReasonKindMismatch)
	requireErrorPath(t, err, "$.name")
}

func TestCompareAddedSubtreeWrapsValueFieldSetUnknownField(t *testing.T) {
	descriptor := types.Object(types.Field("child").Object().Optional()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("child", valueObject("extra", "new")))

	_, err := Compare(valueObject(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.child.extra")
}

func TestCompareRemovedSubtreeWrapsValueFieldSetUnknownField(t *testing.T) {
	descriptor := types.Object(types.Field("child").Object().Optional()).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("child", valueObject("extra", "old")))

	_, err := Compare(oldValue, valueObject(), descriptor, Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.child.extra")
}

func TestCompareAddedSubtreeWrapsInvalidListKey(t *testing.T) {
	descriptor := types.Object(
		types.Field("conditions").ListOf(conditionExpr()).Map("type").Optional(),
	).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember(
		"conditions",
		value.MustListValue(valueObject("status", "True")),
	))

	_, err := Compare(valueObject(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareRemovedSubtreeWrapsInvalidListKey(t *testing.T) {
	descriptor := types.Object(
		types.Field("conditions").ListOf(conditionExpr()).Map("type").Optional(),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember(
		"conditions",
		value.MustListValue(valueObject("status", "True")),
	))

	_, err := Compare(oldValue, valueObject(), descriptor, Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareAddedSubtreeWrapsUnresolvedRef(t *testing.T) {
	descriptor := types.Object(types.Field("name").Ref("example.Name").Optional()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("name", value.StringValue("api")))

	_, err := Compare(valueObject(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
	requireErrorPath(t, err, "$.name")
}

func TestCompareRemovedSubtreeWrapsReferenceCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.Ref("example.A")),
	}
	descriptor := types.Object(types.Field("name").Ref("example.A").Optional()).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("name", value.StringValue("api")))

	_, err := Compare(oldValue, valueObject(), descriptor, Options{Resolver: resolver})

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$.name")
}
