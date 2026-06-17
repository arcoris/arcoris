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

package objectquery

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestErrorHelpersPreserveSentinels verifies representative error builders
// retain both broad and specific errors.Is classifications.
func TestErrorHelpersPreserveSentinels(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
	}{
		{name: "expression", err: invalidExpressionError("bad"), target: ErrInvalidExpression},
		{name: "term", err: invalidTermError("bad"), target: ErrInvalidTerm},
		{name: "field", err: invalidFieldError(ErrInvalidTerm, "bad"), target: ErrInvalidField},
		{name: "operator", err: invalidOperatorError(Operator(99)), target: ErrInvalidOperator},
		{name: "unsupported", err: unsupportedOperatorError(OperatorContains, "metadata"), target: ErrUnsupportedOperator},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, ErrInvalidQuery)
			requireErrorIs(t, tt.err, tt.target)
		})
	}
}

// TestErrorHelpersExposeStructuredDiagnostics verifies callers can recover a
// stable diagnostic path and reason without string parsing.
func TestErrorHelpersExposeStructuredDiagnostics(t *testing.T) {
	err := unresolvedFieldError(fieldRef("spec.image"), ErrUnresolvedField)

	var queryErr *Error
	if !errors.As(err, &queryErr) {
		t.Fatal("errors.As(*Error) = false; want true")
	}
	if queryErr.Reason != ErrorReasonUnresolvedField {
		t.Fatalf("reason = %s; want %s", queryErr.Reason, ErrorReasonUnresolvedField)
	}
	if !strings.Contains(queryErr.Path, "desired.spec.image") {
		t.Fatalf("path = %q; want desired.spec.image", queryErr.Path)
	}
}

// TestCompileUnresolvedFieldErrorPreservesFieldRef verifies field registry
// misses keep both the broad sentinel and a field-specific diagnostic.
func TestCompileUnresolvedFieldErrorPreservesFieldRef(t *testing.T) {
	ref := fieldRef("spec.image")
	_, err := Compile(mustQ(FieldEquals(ref, value.StringValue("api"))), WithSelectableFields(mustFieldSet(t)))

	requireErrorIs(t, err, ErrInvalidQuery)
	requireErrorIs(t, err, ErrUnresolvedField)
	requireStructuredQueryError(t, err, ErrorReasonUnresolvedField, "desired.spec.image")
}

// TestCompileUnsupportedOperatorErrorPreservesDomain verifies operator
// failures identify the invalid field/operator boundary.
func TestCompileUnsupportedOperatorErrorPreservesDomain(t *testing.T) {
	ref := fieldRef("spec.image")
	fields := mustFieldSet(t, selectable(ref, value.KindString, Operators(OperatorEquals)))
	_, err := Compile(mustQ(FieldHasPrefix(ref, "api")), WithSelectableFields(fields))

	requireErrorIs(t, err, ErrUnsupportedOperator)
	requireStructuredQueryError(t, err, ErrorReasonUnsupportedOperator, "query.operator")
}

// TestCompileInvalidMetadataKeyErrorPreservesTermPath verifies metadata
// validation remains distinguishable from field validation.
func TestCompileInvalidMetadataKeyErrorPreservesTermPath(t *testing.T) {
	_, err := Compile(termQuery(term{
		kind:           termMetadata,
		metadataDomain: metadataLabels,
		operator:       OperatorExists,
	}))

	requireErrorIs(t, err, ErrInvalidTerm)
	requireStructuredQueryError(t, err, ErrorReasonInvalidTerm, "query.term")
}

// TestCompileInvalidLiteralKindErrorPreservesFieldDiagnostic verifies literal
// kind mismatches report the field validation boundary.
func TestCompileInvalidLiteralKindErrorPreservesFieldDiagnostic(t *testing.T) {
	ref := fieldRef("spec.replicas")
	fields := mustFieldSet(t, selectable(ref, value.KindInteger, Operators(OperatorEquals)))
	_, err := Compile(mustQ(FieldEquals(ref, value.StringValue("2"))), WithSelectableFields(fields))

	requireErrorIs(t, err, ErrInvalidField)
	requireStructuredQueryError(t, err, ErrorReasonInvalidField, "query.field")
}

func requireStructuredQueryError(t testing.TB, err error, reason ErrorReason, pathPart string) {
	t.Helper()

	var queryErr *Error
	if !errors.As(err, &queryErr) {
		t.Fatal("errors.As(*Error) = false; want true")
	}
	if queryErr.Reason != reason {
		t.Fatalf("reason = %s; want %s", queryErr.Reason, reason)
	}
	if !strings.Contains(queryErr.Path, pathPart) {
		t.Fatalf("path = %q; want to contain %q", queryErr.Path, pathPart)
	}
}
