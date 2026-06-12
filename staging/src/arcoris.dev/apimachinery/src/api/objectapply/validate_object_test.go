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

package objectapply

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/objectvalidation"
	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateObjectMeta(t *testing.T) {
	err := validateObjectMeta(pathObjectLive, testRequest().Live, ErrorReasonInvalidLiveObject)

	requireNoError(t, err)
}

func TestValidateObjectMetaRejectsInvalidMetadata(t *testing.T) {
	obj := testRequest().Live
	obj.TypeMeta = meta.TypeMeta{Kind: "Worker"}

	err := validateObjectMeta(pathObjectLive, obj, ErrorReasonInvalidLiveObject)

	requireErrorIs(t, err, ErrInvalidObject)
}

func TestValidateObject(t *testing.T) {
	req := testRequest()

	err := New(Options{}).validateObject(
		pathObjectLive,
		req.Live,
		ErrorReasonInvalidLiveObject,
		req,
	)

	requireNoError(t, err)
}

func TestValidateObjectMapsVersionMissingAsObjectFailure(t *testing.T) {
	req := testRequest()
	req.Resource = resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion("v2", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	err := New(Options{}).validateObject(
		pathObjectLive,
		req.Live,
		ErrorReasonInvalidLiveObject,
		req,
	)

	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, objectvalidation.ErrVersionNotDefined)
	requireObjectApplyError(t, err, pathObjectLive, ErrorReasonInvalidLiveObject)
}

func TestValidateObjectMapsResourceMismatchAsObjectFailure(t *testing.T) {
	req := testRequest()
	req.Resource = resource.NewDefinition(
		apiidentity.Group("other.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	err := New(Options{}).validateObject(
		pathObjectLive,
		req.Live,
		ErrorReasonInvalidLiveObject,
		req,
	)

	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, objectvalidation.ErrResourceMismatch)
	requireObjectApplyError(t, err, pathObjectLive, ErrorReasonInvalidLiveObject)
}

func TestValidateObjectMapsInvalidPlanAsInvalidResource(t *testing.T) {
	err := objectValidationError(
		pathObjectLive,
		ErrorReasonInvalidLiveObject,
		objectvalidation.ErrInvalidPlan,
	)

	requireErrorIs(t, err, ErrInvalidResource)
	requireErrorIs(t, err, objectvalidation.ErrInvalidPlan)
	requireObjectApplyError(t, err, pathRequestResource, ErrorReasonInvalidResource)
}
