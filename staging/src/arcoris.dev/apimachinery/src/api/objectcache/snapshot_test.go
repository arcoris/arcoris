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

func TestSnapshotIsZeroLenAndRevision(t *testing.T) {
	zero := Snapshot{}
	if !zero.IsZero() {
		t.Fatal("zero Snapshot IsZero() = false; want true")
	}
	if got := zero.Len(); got != 0 {
		t.Fatalf("zero Snapshot Len() = %d; want 0", got)
	}
	if got := zero.Revision(); !got.IsZero() {
		t.Fatalf("zero Snapshot Revision() = %v; want zero", got)
	}

	snapshot := mustSnapshot(t, testListResult(9, testItems()...))
	if snapshot.IsZero() {
		t.Fatal("non-empty Snapshot IsZero() = true; want false")
	}
	if got := snapshot.Len(); got != len(testItems()) {
		t.Fatalf("Len() = %d; want %d", got, len(testItems()))
	}
	if got := snapshot.Revision(); got != 9 {
		t.Fatalf("Revision() = %v; want 9", got)
	}
}
