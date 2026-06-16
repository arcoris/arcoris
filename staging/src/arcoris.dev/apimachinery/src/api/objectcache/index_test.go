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

func TestNewIndexesInitializesBuckets(t *testing.T) {
	idx := newIndexes()

	if idx.byNamespace == nil || idx.byName == nil || idx.byObject == nil {
		t.Fatalf("identity buckets not initialized: %#v", idx)
	}
	if idx.byLabelKey == nil || idx.byLabelValue == nil {
		t.Fatalf("label buckets not initialized: %#v", idx)
	}
	if idx.byAnnotationKey == nil || idx.byAnnotationValue == nil {
		t.Fatalf("annotation buckets not initialized: %#v", idx)
	}
}

func TestIndexKeyBucketAddRemove(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil)
	buckets := map[string]keySet{}

	addIndexKey(buckets, "bucket", item.Key)
	if !buckets["bucket"].has(item.Key) {
		t.Fatal("bucket missing key after add")
	}

	removeIndexKey(buckets, "bucket", item.Key)
	if _, exists := buckets["bucket"]; exists {
		t.Fatal("empty bucket still exists after remove")
	}
}
