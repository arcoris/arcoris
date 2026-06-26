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

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCollectionRemoveOrderKeyPreservesRelativeOrder(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	third := testKey("system", 3)
	col := collection{order: []objectstore.Key{first, second, third}}

	col.removeOrderKey(second)

	if len(col.order) != 2 || !col.order[0].Equal(first) || !col.order[1].Equal(third) {
		t.Fatalf("order = %#v; want first, third", col.order)
	}
}
