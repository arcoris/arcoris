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

package fieldownership

import "testing"

func TestOwnedPathValidateStructureAcceptsValidRecord(t *testing.T) {
	err := OwnedPath{Owner: owner("user-cli"), Path: imagePath()}.ValidateStructure()

	requireNoError(t, err)
}

func TestOwnedPathValidateStructureAcceptsRootPath(t *testing.T) {
	err := OwnedPath{Owner: owner("user-cli")}.ValidateStructure()

	requireNoError(t, err)
}

func TestOwnedPathValidateStructureRejectsInvalidOwner(t *testing.T) {
	err := OwnedPath{Owner: Owner{}, Path: imagePath()}.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidOwnedPath)
	requireErrorIs(t, err, ErrInvalidOwner)
	requireFieldOwnershipError(t, err, "ownedPath.owner", ErrorReasonInvalidOwnedPathOwner)
}
