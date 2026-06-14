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

	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestApplyDesiredAbsentOwnedFieldDeletesIt(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("user", path("$.image"), path("$.replicas")))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "image", "api:v2")
	requireNoMember(t, result.Object.Desired, "replicas")
	requireSet(t, result.Desired.DeletedFields, "$.replicas")
	requireOwnersOf(t, result.Ownership.Desired(), path("$.image"), "user")
	owners, err := result.Ownership.Desired().OwnersOf(path("$.replicas"))
	requireNoError(t, err)
	requireOwners(t, owners)
}

func TestApplyDesiredExplicitNullUsesDescriptorNullability(t *testing.T) {
	req := testRequest()
	req.Applied = appliedObject(obj(member("image", value.NullValue())))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, valuevalidation.ErrNullNotAllowed)
	requireObjectApplyError(t, err, pathObjectApplied, ErrorReasonInvalidAppliedObject)
}
