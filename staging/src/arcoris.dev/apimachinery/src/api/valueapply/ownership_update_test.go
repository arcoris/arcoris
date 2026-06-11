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
	merged := mergedApply{
		preparedApply: preparedApply{AppliedFields: fields(imagePath())},
		Value:         req.Live,
	}

	got, err := newApplier(Options{}).updateOwnership(req, merged)
	requireNoError(t, err)

	requireSet(t, got.Ownership.FieldsFor(owner("user")), "$.image")
}

func TestUpdateOwnershipForceRemovesConflictingOthers(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))
	merged := mergedApply{
		preparedApply: preparedApply{
			AppliedFields: fields(imagePath()),
			Conflicts:     mustConflicts(req, fields(imagePath())),
		},
		Value: req.Live,
	}

	got, err := newApplier(Options{Force: true}).updateOwnership(req, merged)
	requireNoError(t, err)

	requireOwnersOf(t, got.Ownership, imagePath(), "user")
}

func TestUpdateOwnershipWrapsInvalidOwner(t *testing.T) {
	req := specRequest(owner("user"))
	req.Owner = fieldownership.Owner{}
	merged := mergedApply{
		preparedApply: preparedApply{AppliedFields: fields(imagePath())},
		Value:         req.Live,
	}

	_, err := newApplier(Options{}).updateOwnership(req, merged)

	requireErrorIs(t, err, ErrInvalidRequest)
}

func mustConflicts(req Request, attempted fieldpath.Set) fieldownership.ConflictSet {
	conflicts, err := req.Ownership.Conflicts(req.Owner, attempted)
	if err != nil {
		panic(err)
	}

	return conflicts
}
