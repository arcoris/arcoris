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

func TestBuildCollectionEmptyNormalizesMaps(t *testing.T) {
	col := mustCollection(t, testListResult(5))

	if col.items != nil {
		t.Fatalf("items = %#v; want nil", col.items)
	}
	if col.indexes.byNamespace != nil {
		t.Fatalf("indexes = %#v; want zero indexes", col.indexes)
	}
	if col.revision != 5 {
		t.Fatalf("revision = %v; want 5", col.revision)
	}
}

func TestBuildCollectionDuplicateKeyError(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil)

	_, err := buildCollection(testListResult(2, item, item), ErrInvalidSnapshot)

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, ErrDuplicateKey)
}
