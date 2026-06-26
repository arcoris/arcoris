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

func TestBuildCollectionPreservesOrderRevisionAndDetachesItems(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	result := objectstore.ListResult{
		Items: []objectstore.ListItem{
			listItem(first, 1, "first"),
			listItem(second, 2, "second"),
		},
		Revision: 3,
	}

	col, err := buildCollection(result)
	requireNoError(t, err)
	mutateState(&result.Items[0].State, "mutated")

	got := col.listResult()
	requireListOrder(t, got, first, second)
	if got.Revision != 3 {
		t.Fatalf("revision = %s; want 3", got.Revision)
	}
	if desired := desiredString(t, got.Items[0].State); desired != "first" {
		t.Fatalf("desired = %q; want first", desired)
	}
}

func TestBuildCollectionRejectsDuplicateKeys(t *testing.T) {
	key := testKey("system", 1)
	_, err := buildCollection(objectstore.ListResult{
		Items: []objectstore.ListItem{
			listItem(key, 1, "first"),
			listItem(key, 2, "second"),
		},
		Revision: 2,
	})

	requireErrorIs(t, err, ErrInvalidRead)
	requireErrorIs(t, err, ErrDuplicateKey)
}
