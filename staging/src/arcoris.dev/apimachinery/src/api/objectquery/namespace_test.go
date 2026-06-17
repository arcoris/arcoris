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

// TestPredicateNamespaceRequirement verifies namespace extraction is explicit
// and conservative enough for lifecycle pre-validation.
func TestPredicateNamespaceRequirement(t *testing.T) {
	alpha := mustQ(ObjectInNamespace("alpha"))
	beta := mustQ(ObjectInNamespace("beta"))
	alphaObject := mustQ(ObjectEquals("alpha", "worker"))
	alphaKey := mustQ(KeyEquals(testItem("alpha", "worker", 1, nil, nil, desiredRecord("api", "prod", 1)).Key))

	tests := []struct {
		name string
		q    Query
		kind NamespaceRequirementKind
		ns   metaidentity.Namespace
	}{
		{name: "all", q: All(), kind: NamespaceUnconstrained},
		{name: "namespace", q: alpha, kind: NamespaceSingle, ns: "alpha"},
		{name: "object", q: alphaObject, kind: NamespaceSingle, ns: "alpha"},
		{name: "key", q: alphaKey, kind: NamespaceSingle, ns: "alpha"},
		{name: "and same", q: mustAnd(t, alpha, alphaObject), kind: NamespaceSingle, ns: "alpha"},
		{name: "and conflict", q: mustAnd(t, alpha, beta), kind: NamespaceContradictory},
		{name: "or same", q: mustOr(t, alpha, alphaObject), kind: NamespaceSingle, ns: "alpha"},
		{name: "or different", q: mustOr(t, alpha, beta), kind: NamespaceDisjunctive},
		{name: "not", q: mustNot(t, alpha), kind: NamespaceResidual},
		{name: "or unconstrained", q: mustOr(t, alpha, mustQ(LabelEquals("env", "prod"))), kind: NamespaceUnconstrained},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustPredicate(t, tt.q).NamespaceRequirement()
			if got.Kind != tt.kind || got.Namespace != tt.ns {
				t.Fatalf("NamespaceRequirement = (%v, %q); want (%v, %q)", got.Kind, got.Namespace, tt.kind, tt.ns)
			}
		})
	}
}
