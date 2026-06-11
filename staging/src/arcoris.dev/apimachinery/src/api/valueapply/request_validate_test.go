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

package valueapply

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

func TestValidateRequestShape(t *testing.T) {
	requireNoError(t, validateRequestShape(specRequest(owner("user"))))
}

func TestValidateRequestShapeInvalidPath(t *testing.T) {
	req := specRequest(owner("user"))
	req.Path = root().Append(fieldpath.Element{})

	err := validateRequestShape(req)

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, fieldpath.ErrInvalidPath)
}

func TestValidateRequestShapeOwnershipEqualToBaseAccepted(t *testing.T) {
	req := specRequest(owner("user"))
	req.Path = specPath()
	req.Ownership = state(entry("user", specPath()))

	requireNoError(t, validateRequestShape(req))
}

func TestValidateRequestShapeOwnershipDescendantOfBaseAccepted(t *testing.T) {
	req := specRequest(owner("user"))
	req.Path = specPath()
	req.Ownership = state(entry("user", specPath().Field(testFieldName("image"))))

	requireNoError(t, validateRequestShape(req))
}

func TestValidateRequestShapeOwnershipOutsideBaseRejected(t *testing.T) {
	req := specRequest(owner("user"))
	req.Path = specPath()
	req.Ownership = state(entry("user", path("$.metadata.name")))

	err := validateRequestShape(req)

	requireErrorIs(t, err, ErrInvalidRequest)
}

func TestApplyOwnershipOutsideBaseReturnsInvalidRequest(t *testing.T) {
	req := specRequest(owner("user"))
	req.Path = specPath()
	req.Live = obj(member("image", str("api:v1")))
	req.Applied = obj(member("image", str("api:v2")))
	req.Descriptor = specStringDescriptor()
	req.Ownership = state(entry("user", path("$.metadata.name")))

	result, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidRequest)
	if !result.IsZero() {
		t.Fatalf("result is not zero")
	}
}

func TestApplyInvalidOwnerReturnsInvalidOwner(t *testing.T) {
	_, err := Apply(Request{
		Path:       root(),
		Owner:      fieldownership.Owner{},
		Live:       str("old"),
		Applied:    str("new"),
		Descriptor: types.String().Descriptor(),
	}, Options{})

	requireErrorIs(t, err, ErrInvalidOwner)
}
