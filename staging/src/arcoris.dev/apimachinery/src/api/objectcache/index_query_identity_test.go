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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
)

func TestIndexesPlanIdentityNamespaceNameAndObject(t *testing.T) {
	idx := testIndexes()
	tests := []struct {
		name     string
		identity objectquery.IdentitySelector
		want     []itemRef
	}{
		{
			name: "namespace",
			identity: objectquery.IdentitySelector{
				Namespace: mustNamespaceEquals(t, "system"),
			},
			want: []itemRef{
				{"system", "worker-1", 1},
				{"system", "worker-2", 2},
				{"system", "worker-4", 4},
			},
		},
		{
			name: "name",
			identity: objectquery.IdentitySelector{
				Name: mustNameEquals(t, "worker-3"),
			},
			want: []itemRef{{"other", "worker-3", 3}},
		},
		{
			name: "object",
			identity: objectquery.IdentitySelector{
				Namespace: mustNamespaceEquals(t, "system"),
				Name:      mustNameEquals(t, "worker-2"),
			},
			want: []itemRef{{"system", "worker-2", 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := idx.planIdentity(candidatePlan{}, tt.identity)

			requirePlanIncludesOnly(t, plan, tt.want...)
		})
	}
}

func requirePlanIncludesOnly(t *testing.T, plan candidatePlan, want ...itemRef) {
	t.Helper()

	expected := map[itemRef]struct{}{}
	for _, ref := range want {
		expected[ref] = struct{}{}
	}
	for _, item := range testItems() {
		ref := itemRef{
			namespace: item.Key.Object.Namespace,
			name:      item.Key.Object.Name,
			revision:  item.State.Revision,
		}
		_, shouldInclude := expected[ref]
		if got := plan.includes(item.Key); got != shouldInclude {
			t.Fatalf("plan.includes(%s) = %t; want %t", item.Key.String(), got, shouldInclude)
		}
	}
}
