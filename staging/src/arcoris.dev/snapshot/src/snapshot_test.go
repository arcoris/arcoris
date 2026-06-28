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

package snapshot

import "testing"

type testRevision string

func TestSnapshotZeroValue(t *testing.T) {
	var snap Snapshot[LocalRevision, string]

	if !snap.Revision.IsZero() {
		t.Fatalf("zero Snapshot revision = %d, want zero", snap.Revision)
	}
	if got, want := snap.Value, ""; got != want {
		t.Fatalf("zero Snapshot value = %q, want %q", got, want)
	}
	if !snap.ChangedSince(LocalRevision(1)) {
		t.Fatal("zero Snapshot should differ from a committed revision")
	}
}

func TestSnapshotRevisionHelpers(t *testing.T) {
	snap := Snapshot[LocalRevision, string]{Revision: LocalRevision(2), Value: "value"}

	if snap.Revision.IsZero() {
		t.Fatal("non-zero snapshot revision reported zero")
	}

	if snap.ChangedSince(LocalRevision(2)) {
		t.Fatal("snapshot changed since same revision")
	}

	if !snap.ChangedSince(LocalRevision(1)) {
		t.Fatal("snapshot did not change since different revision")
	}
}

func TestSnapshotCarriesCustomRevisionType(t *testing.T) {
	snap := Snapshot[testRevision, string]{
		Revision: testRevision("domain-r2"),
		Value:    "value",
	}

	if got, want := snap.Revision, testRevision("domain-r2"); got != want {
		t.Fatalf("Revision = %q, want %q", got, want)
	}
	if snap.ChangedSince(testRevision("domain-r2")) {
		t.Fatal("custom revision changed since same revision")
	}
	if !snap.ChangedSince(testRevision("domain-r1")) {
		t.Fatal("custom revision did not change since different revision")
	}
}

func TestSnapshotWithValuePreservesRevision(t *testing.T) {
	snap := Snapshot[LocalRevision, string]{Revision: LocalRevision(7), Value: "old"}
	got := snap.WithValue("new")

	if got.Revision != snap.Revision {
		t.Fatalf("WithValue revision = %d, want %d", got.Revision, snap.Revision)
	}
	if got.Value != "new" {
		t.Fatalf("WithValue value = %q, want %q", got.Value, "new")
	}
}
