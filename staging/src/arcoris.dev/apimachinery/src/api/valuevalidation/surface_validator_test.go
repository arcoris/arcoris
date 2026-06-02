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
	"errors"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectvalidation"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestSurfaceValidatorValidatesValue(t *testing.T) {
	adapter := valuevalidation.SurfaceValidator{}

	err := adapter.ValidateSurface(
		value.StringValue(""),
		types.String().MinLen(1).Type(),
		nil,
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$",
	)
}

func TestSurfaceValidatorUsesPlanResolverArgument(t *testing.T) {
	optionsResolver := testResolver{
		"example.Name": types.Define("example.Name", types.String().MinLen(4)),
	}
	planResolver := testResolver{
		"example.Name": types.Define("example.Name", types.String().MinLen(1)),
	}
	adapter := valuevalidation.SurfaceValidator{
		Options: valuevalidation.Options{Resolver: optionsResolver},
	}

	requireNoError(
		t,
		adapter.ValidateSurface(
			value.StringValue("ok"),
			types.Ref("example.Name").Type(),
			planResolver,
		),
	)
}

func TestSurfaceValidatorFallsBackToOptionsResolverWhenArgumentNil(t *testing.T) {
	optionsResolver := testResolver{
		"example.Name": types.Define("example.Name", types.String().MinLen(1)),
	}
	adapter := valuevalidation.SurfaceValidator{
		Options: valuevalidation.Options{Resolver: optionsResolver},
	}

	requireNoError(
		t,
		adapter.ValidateSurface(
			value.StringValue("ok"),
			types.Ref("example.Name").Type(),
			nil,
		),
	)
}

func TestObjectValidationUsesValueSurfaceValidator(t *testing.T) {
	desired := types.Object(
		types.Field("replicas").Int32().Required(),
	).Type()
	resourceDefinition := resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion(apiidentity.Version("v1"), desired),
	)
	payload := mustObject(t, value.ObjectMember("replicas", value.StringValue("three")))
	obj := object.New[value.Value, value.Value](
		meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
			Group:   "control.arcoris.dev",
			Version: "v1",
			Kind:    "Worker",
		}),
		meta.ObjectMeta{
			Name:      metaidentity.Name("worker"),
			Namespace: metaidentity.Namespace("system"),
		},
		payload,
	)
	plan := objectvalidation.Plan[value.Value, value.Value]{
		Resource:         resourceDefinition,
		DesiredValidator: valuevalidation.SurfaceValidator{},
	}

	err := objectvalidation.Validate(obj, plan)
	if !errors.Is(err, objectvalidation.ErrInvalidDesired) {
		t.Fatalf("errors.Is(ErrInvalidDesired) = false: %v", err)
	}
	if !errors.Is(err, valuevalidation.ErrKindMismatch) {
		t.Fatalf("errors.Is(valuevalidation.ErrKindMismatch) = false: %v", err)
	}
}
