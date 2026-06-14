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

func TestIdentitySelectorZeroMatchesEverything(t *testing.T) {
	var selector IdentitySelector

	if !selector.IsZero() {
		t.Fatal("zero identity selector is not zero")
	}
	if !selector.match(testItem("system", "worker", nil, nil)) {
		t.Fatal("zero identity selector did not match")
	}
	requireNoError(t, selector.Validate())
}

func mustInNamespace(t *testing.T, namespace metaidentity.Namespace) IdentitySelector {
	t.Helper()
	selector, err := InNamespace(namespace)
	requireNoError(t, err)

	return selector
}

func mustWithName(t *testing.T, name metaidentity.Name) IdentitySelector {
	t.Helper()
	selector, err := WithName(name)
	requireNoError(t, err)

	return selector
}

func mustWithObject(t *testing.T, namespace metaidentity.Namespace, name metaidentity.Name) IdentitySelector {
	t.Helper()
	selector, err := WithObject(namespace, name)
	requireNoError(t, err)

	return selector
}
