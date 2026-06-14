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

package objectquery

import (
	"testing"

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestIdentitySelectorRejectsInvalidInputs(t *testing.T) {
	_, err := InNamespace(metaidentity.Namespace("Bad_Namespace"))
	requireErrorIs(t, err, ErrInvalidQuery)
	requireErrorIs(t, err, metaidentity.ErrInvalidNamespace)

	_, err = WithName(metaidentity.Name(""))
	requireErrorIs(t, err, ErrInvalidQuery)
	requireErrorIs(t, err, metaidentity.ErrInvalidName)
}
