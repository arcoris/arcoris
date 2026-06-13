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

package objectstore

import (
	"testing"

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestListScopeConstructors(t *testing.T) {
	all := AllNamespaces()
	if all.IsZero() || !all.IsAllNamespaces() || all.IsNamespace() {
		t.Fatalf("AllNamespaces() = %#v; want all non-zero scope", all)
	}
	if all.Kind() != ListScopeAll {
		t.Fatalf("all kind = %s; want all", all.Kind())
	}
	if !all.Namespace().IsZero() {
		t.Fatalf("all namespace = %q; want empty", all.Namespace())
	}

	namespace, err := InNamespace("system")
	requireNoError(t, err)
	if namespace.IsZero() || !namespace.IsNamespace() || namespace.IsAllNamespaces() {
		t.Fatalf("InNamespace() = %#v; want namespace scope", namespace)
	}
	if namespace.Kind() != ListScopeNamespace {
		t.Fatalf("namespace kind = %s; want namespace", namespace.Kind())
	}
	if namespace.Namespace() != metaidentity.Namespace("system") {
		t.Fatalf("namespace = %q; want system", namespace.Namespace())
	}
}

func TestMustNamespacePanicsOnInvalidNamespace(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustNamespace did not panic")
		}
	}()

	_ = MustNamespace("System")
}
