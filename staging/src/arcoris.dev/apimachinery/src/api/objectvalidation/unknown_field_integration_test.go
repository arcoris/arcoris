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

package objectvalidation

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	apiobject "arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateUnknownDesiredFieldCausePreserved(t *testing.T) {
	plan := valueValidationPlan()
	obj := valueObject(
		valueRecord(valueField("unexpected", value.StringValue("x"))),
		nil,
	)

	err := Validate(obj, plan)
	requireValidationError(t, err, ErrInvalidDesired, pathObjectDesired, ErrorReasonInvalidDesired)
	requireErrorIs(t, err, valuevalidation.ErrUnknownField)
}

func TestValidateUnknownObservedFieldCausePreserved(t *testing.T) {
	plan := valueValidationPlan()
	observed := valueRecord(valueField("unexpected", value.StringValue("x")))
	obj := valueObject(
		valueRecord(valueField("image", value.StringValue("api:v1"))),
		&observed,
	)

	err := Validate(obj, plan)
	requireValidationError(t, err, ErrInvalidObserved, pathObjectObserved, ErrorReasonInvalidObserved)
	requireErrorIs(t, err, valuevalidation.ErrUnknownField)
}

func valueValidationPlan() Plan[value.Value, value.Value] {
	validator := valuevalidation.SurfaceValidator{}

	return Plan[value.Value, value.Value]{
		Resource: resource.NewDefinition(
			apiidentity.Group("control.arcoris.dev"),
			apiidentity.Kind("Worker"),
			apiidentity.Resource("workers"),
			resource.ScopeNamespaced,
			resource.NewVersion("v1", valueDesiredDescriptor(), resource.Observed(valueObservedDescriptor())),
		),
		DesiredValidator:  validator,
		ObservedValidator: validator,
	}
}

func valueDesiredDescriptor() types.Descriptor {
	return types.Object(
		types.Field("image").String().Optional(),
	).UnknownFields(types.UnknownReject).Descriptor()
}

func valueObservedDescriptor() types.Descriptor {
	return types.Object(
		types.Field("ready").String().Optional(),
	).UnknownFields(types.UnknownReject).Descriptor()
}

func valueObject(desired value.Value, observed *value.Value) apiobject.Object[value.Value, value.Value] {
	obj := apiobject.New[value.Value, value.Value](
		meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
			Group:   "control.arcoris.dev",
			Version: "v1",
			Kind:    "Worker",
		}),
		meta.ObjectMeta{
			Name:      metaidentity.Name("worker"),
			Namespace: metaidentity.Namespace("system"),
		},
		desired,
	)
	if observed != nil {
		obj = obj.WithObserved(*observed)
	}

	return obj
}

func valueRecord(members ...value.RecordMember) value.Value {
	return value.MustRecordValue(members...)
}

func valueField(name string, val value.Value) value.RecordMember {
	return value.MustRecordMember(name, val)
}
