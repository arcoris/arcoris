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

import "testing"

func TestIdentitySelectorMatchesNamespaceAndName(t *testing.T) {
	tests := []struct {
		name     string
		selector IdentitySelector
		item     objectName
		match    bool
	}{
		{
			name:     "namespace",
			selector: mustInNamespace(t, "system"),
			item:     objectName{namespace: "system", name: "worker"},
			match:    true,
		},
		{
			name:     "namespace mismatch",
			selector: mustInNamespace(t, "system"),
			item:     objectName{namespace: "other", name: "worker"},
		},
		{
			name:     "name",
			selector: mustWithName(t, "worker"),
			item:     objectName{namespace: "system", name: "worker"},
			match:    true,
		},
		{
			name:     "name mismatch",
			selector: mustWithName(t, "worker"),
			item:     objectName{namespace: "system", name: "other"},
		},
		{
			name:     "object",
			selector: mustWithObject(t, "system", "worker"),
			item:     objectName{namespace: "system", name: "worker"},
			match:    true,
		},
		{
			name:     "object namespace mismatch",
			selector: mustWithObject(t, "system", "worker"),
			item:     objectName{namespace: "other", name: "worker"},
		},
		{
			name:     "empty namespace exact",
			selector: mustInNamespace(t, ""),
			item:     objectName{name: "worker"},
			match:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := testItem(tt.item.namespace, tt.item.name, nil, nil)
			if got := tt.selector.match(item); got != tt.match {
				t.Fatalf("match = %v; want %v", got, tt.match)
			}
		})
	}
}

type objectName struct {
	namespace string
	name      string
}
