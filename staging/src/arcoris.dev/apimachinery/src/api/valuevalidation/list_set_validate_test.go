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

package valuevalidation_test

import (
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateListSetAcceptsDistinctStableScalarElements(t *testing.T) {
	tests := []struct {
		name       string
		descriptor types.Descriptor
		payload    value.Value
	}{
		{
			name:       "strings",
			descriptor: types.ListOf(types.String()).Set().Descriptor(),
			payload:    mustList(t, value.StringValue("a"), value.StringValue("b"), value.StringValue("c")),
		},
		{
			name:       "signed integers",
			descriptor: types.ListOf(types.Int64()).Set().Descriptor(),
			payload:    mustList(t, value.Int64Value(1), value.Int64Value(2), value.Int64Value(3)),
		},
		{
			name:       "unsigned integers",
			descriptor: types.ListOf(types.Uint64()).Set().Descriptor(),
			payload:    mustList(t, value.Uint64Value(1), value.Uint64Value(2), value.Uint64Value(3)),
		},
		{
			name:       "booleans",
			descriptor: types.ListOf(types.Bool()).Set().Descriptor(),
			payload:    mustList(t, value.BoolValue(true), value.BoolValue(false)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireNoError(
				t,
				valuevalidation.Validate(
					tt.payload,
					tt.descriptor,
					valuevalidation.Options{},
				),
			)
		})
	}
}

func TestValidateListSetAcceptsDistinctRefStringElements(t *testing.T) {
	resolver := testResolver{
		"meta.arcoris.dev.Name": types.Define("meta.arcoris.dev.Name", types.String().MinBytes(1)),
	}
	descriptor := types.ListOf(types.Ref("meta.arcoris.dev.Name")).Set().Descriptor()
	payload := mustList(t, value.StringValue("api"), value.StringValue("worker"))

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			descriptor,
			valuevalidation.Options{Resolver: resolver},
		),
	)
}

func TestValidateListSetRejectsDuplicateStringElements(t *testing.T) {
	descriptor := types.ListOf(types.String()).Set().Descriptor()
	payload := mustList(t, value.StringValue("a"), value.StringValue("b"), value.StringValue("a"))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{})

	requireError(
		t,
		err,
		valuevalidation.ErrDuplicateListSetElement,
		valuevalidation.ErrorReasonDuplicateListSetElement,
		"$[2]",
	)

	validationError := findValidationError(
		t,
		err,
		valuevalidation.ErrDuplicateListSetElement,
		valuevalidation.ErrorReasonDuplicateListSetElement,
		"$[2]",
	)
	if !strings.Contains(validationError.Detail, "index 0") {
		t.Fatalf("duplicate detail = %q, want first index", validationError.Detail)
	}
}

func TestValidateListSetRejectsDuplicateSignedIntegerElements(t *testing.T) {
	descriptor := types.ListOf(types.Int64()).Set().Descriptor()
	payload := mustList(t, value.Int64Value(1), value.Int64Value(2), value.Int64Value(1))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{})

	requireError(
		t,
		err,
		valuevalidation.ErrDuplicateListSetElement,
		valuevalidation.ErrorReasonDuplicateListSetElement,
		"$[2]",
	)
}

func TestValidateListSetRejectsDuplicateUnsignedIntegerElements(t *testing.T) {
	descriptor := types.ListOf(types.Uint64()).Set().Descriptor()
	payload := mustList(t, value.Uint64Value(1), value.Uint64Value(2), value.Uint64Value(1))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{})

	requireError(
		t,
		err,
		valuevalidation.ErrDuplicateListSetElement,
		valuevalidation.ErrorReasonDuplicateListSetElement,
		"$[2]",
	)
}

func TestValidateListSetRejectsDuplicateBoolElements(t *testing.T) {
	descriptor := types.ListOf(types.Bool()).Set().Descriptor()
	payload := mustList(t, value.BoolValue(true), value.BoolValue(false), value.BoolValue(true))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{})

	requireError(
		t,
		err,
		valuevalidation.ErrDuplicateListSetElement,
		valuevalidation.ErrorReasonDuplicateListSetElement,
		"$[2]",
	)
}

func TestValidateListSetRejectsDuplicateRefStringElements(t *testing.T) {
	resolver := testResolver{
		"meta.arcoris.dev.Name": types.Define("meta.arcoris.dev.Name", types.String().MinBytes(1)),
	}
	descriptor := types.ListOf(types.Ref("meta.arcoris.dev.Name")).Set().Descriptor()
	payload := mustList(t, value.StringValue("api"), value.StringValue("worker"), value.StringValue("api"))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{Resolver: resolver})

	requireError(
		t,
		err,
		valuevalidation.ErrDuplicateListSetElement,
		valuevalidation.ErrorReasonDuplicateListSetElement,
		"$[2]",
	)
}

func TestValidateListSetRejectsObjectElementDescriptor(t *testing.T) {
	descriptor := types.ListOf(types.Object(
		types.Field("name").String().Required(),
	)).Set().Descriptor()
	payload := mustList(t, mustObject(t, value.MustRecordMember("name", value.StringValue("api"))))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{})

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidDescriptor,
		valuevalidation.ErrorReasonInvalidDescriptor,
		"$",
	)
}

func TestValidateListSetRejectsRefResolvingToObject(t *testing.T) {
	resolver := testResolver{
		"meta.arcoris.dev.Condition": types.Define(
			"meta.arcoris.dev.Condition",
			types.Object(types.Field("type").String().Required()),
		),
	}
	descriptor := types.ListOf(types.Ref("meta.arcoris.dev.Condition")).Set().Descriptor()
	payload := mustList(t, mustObject(t, value.MustRecordMember("type", value.StringValue("Ready"))))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{Resolver: resolver})

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidDescriptor,
		valuevalidation.ErrorReasonInvalidDescriptor,
		"$",
	)
}

func TestValidateListSetRejectsUnresolvedRef(t *testing.T) {
	descriptor := types.ListOf(types.Ref("meta.arcoris.dev.Name")).Set().Descriptor()
	payload := mustList(t, value.StringValue("api"))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{})

	requireError(
		t,
		err,
		valuevalidation.ErrUnresolvedRef,
		valuevalidation.ErrorReasonUnresolvedRef,
		"$",
	)
}

func TestValidateListSetRejectsReferenceCycle(t *testing.T) {
	resolver := testResolver{
		"example.dev.A": types.Define("example.dev.A", types.Ref("example.dev.B")),
		"example.dev.B": types.Define("example.dev.B", types.Ref("example.dev.A")),
	}
	descriptor := types.ListOf(types.Ref("example.dev.A")).Set().Descriptor()
	payload := mustList(t, value.StringValue("api"))

	err := valuevalidation.Validate(payload, descriptor, valuevalidation.Options{Resolver: resolver})

	requireError(
		t,
		err,
		valuevalidation.ErrReferenceCycle,
		valuevalidation.ErrorReasonReferenceCycle,
		"$",
	)
}
