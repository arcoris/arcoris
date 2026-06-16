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

import "testing"

func requireKeySetIncludesOnly(t *testing.T, set keySet, want ...itemRef) {
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
		if got := set.has(item.Key); got != shouldInclude {
			t.Fatalf("set.has(%s) = %t; want %t", item.Key.String(), got, shouldInclude)
		}
	}
}
