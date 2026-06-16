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

func TestIndexesAddAndRemoveItem(t *testing.T) {
	item := testItem(
		"system",
		"worker-1",
		1,
		labelsMap("env", "prod"),
		annotationsMap("team", "core"),
	)
	idx := newIndexes()

	idx.add(item)
	if !idx.byNamespace[item.Key.Object.Namespace].has(item.Key) {
		t.Fatal("namespace bucket missing key after add")
	}
	if !idx.byLabelKey["env"].has(item.Key) {
		t.Fatal("label key bucket missing key after add")
	}
	if !idx.byAnnotationKey["team"].has(item.Key) {
		t.Fatal("annotation key bucket missing key after add")
	}

	idx.remove(item)
	if len(idx.byNamespace) != 0 {
		t.Fatalf("namespace buckets = %#v; want empty", idx.byNamespace)
	}
	if len(idx.byLabelKey) != 0 || len(idx.byLabelValue) != 0 {
		t.Fatalf("label buckets not empty: %#v %#v", idx.byLabelKey, idx.byLabelValue)
	}
	if len(idx.byAnnotationKey) != 0 || len(idx.byAnnotationValue) != 0 {
		t.Fatalf("annotation buckets not empty: %#v %#v", idx.byAnnotationKey, idx.byAnnotationValue)
	}
}
