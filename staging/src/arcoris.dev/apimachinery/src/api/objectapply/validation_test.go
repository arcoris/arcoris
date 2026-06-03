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
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestApplyInvalidOwner(t *testing.T) {
	req := testRequest()
	req.Owner = owner(" ")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidOwner)
}

func TestApplyMissingResource(t *testing.T) {
	req := testRequest()
	req.Resource = resource.Definition{}

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidResource)
}

func TestApplyInvalidLiveObject(t *testing.T) {
	req := testRequest()
	req.Live.TypeMeta = meta.TypeMeta{Kind: "Worker"}

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
}

func TestApplyInvalidAppliedObject(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Name = "Worker"

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
}

func TestApplyResourceMismatch(t *testing.T) {
	req := testRequest()
	req.Resource = resource.NewDefinition(
		apiidentity.Group("other.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
}

func TestApplyScopeMismatch(t *testing.T) {
	req := testRequest()
	req.Resource = resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeGlobal,
		resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
}

func TestApplyDesiredValidationFailure(t *testing.T) {
	req := testRequest()
	req.Applied = appliedObject(obj(member("image", value.Int64Value(1))))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
}

func TestApplyInvalidDesiredValueRejectedByValueApply(t *testing.T) {
	req := testRequest()
	req.Live = testObject(value.Value{})
	req.Applied = appliedObject(str("new"))
	req.Resource = testResource(types.String().Type())

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
}
