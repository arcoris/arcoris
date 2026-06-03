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

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestValidateIdentityCompatibility(t *testing.T) {
	req := testRequest()

	err := validateIdentityCompatibility(req.Live, req.Applied)

	requireNoError(t, err)
}

func TestValidateIdentityCompatibilityRejectsNameMismatch(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Name = metaidentity.Name("other")

	err := validateIdentityCompatibility(req.Live, req.Applied)

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestValidateIdentityCompatibilityRejectsUIDMismatch(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.UID = metaidentity.UID("uid-2")

	err := validateIdentityCompatibility(req.Live, req.Applied)

	requireErrorIs(t, err, ErrIdentityMismatch)
}
