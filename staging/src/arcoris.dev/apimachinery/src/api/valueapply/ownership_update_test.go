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
)

func TestUpdateOwnershipSetsOwnerFields(t *testing.T) {
	req := specRequest(owner("user"))
	result := Result{AppliedFields: fields(imagePath())}

	got, err := newApplier(Options{}).updateOwnership(req, result)
	requireNoError(t, err)

	requireSet(t, got.FieldsFor(owner("user")), "$.image")
}

func TestUpdateOwnershipForceRemovesConflictingOthers(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))
	result := Result{
		AppliedFields: fields(imagePath()),
		Conflicts:     mustConflicts(req, fields(imagePath())),
	}

	got, err := newApplier(Options{Force: true}).updateOwnership(req, result)
	requireNoError(t, err)

	requireOwners(t, got.OwnersOf(imagePath()), "user")
}

func TestUpdateOwnershipWrapsInvalidOwner(t *testing.T) {
	req := specRequest(owner(" "))
	result := Result{AppliedFields: fields(imagePath())}

	_, err := newApplier(Options{}).updateOwnership(req, result)

	requireErrorIs(t, err, ErrInvalidRequest)
}

func mustConflicts(req Request, attempted fieldpath.Set) fieldownership.ConflictSet {
	conflicts, err := req.Ownership.Conflicts(req.Owner, attempted)
	if err != nil {
		panic(err)
	}

	return conflicts
}
