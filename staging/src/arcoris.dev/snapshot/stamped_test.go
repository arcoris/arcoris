/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package snapshot

import (
	"testing"
	"time"
)

func TestStampedRevisionHelpers(t *testing.T) {
	stamped := Stamped[string]{Revision: Revision(3), Updated: time.Unix(1, 0), Value: "value"}

	if stamped.IsZeroRevision() {
		t.Fatal("non-zero stamped revision reported zero")
	}

	if stamped.ChangedSince(Revision(3)) {
		t.Fatal("stamped changed since same revision")
	}

	if !stamped.ChangedSince(Revision(2)) {
		t.Fatal("stamped did not change since different revision")
	}
}

func TestStampedAge(t *testing.T) {
	updated := time.Unix(10, 0)
	stamped := Stamped[string]{Revision: Revision(1), Updated: updated, Value: "value"}

	if got, want := stamped.Age(updated.Add(5*time.Second)), 5*time.Second; got != want {
		t.Fatalf("Age() = %s, want %s", got, want)
	}
}

func TestStampedWithValuePreservesMetadata(t *testing.T) {
	updated := time.Unix(10, 0)
	stamped := Stamped[string]{Revision: Revision(8), Updated: updated, Value: "old"}
	got := stamped.WithValue("new")

	if got.Revision != stamped.Revision {
		t.Fatalf("WithValue revision = %d, want %d", got.Revision, stamped.Revision)
	}
	if !got.Updated.Equal(updated) {
		t.Fatalf("WithValue updated = %s, want %s", got.Updated, updated)
	}
	if got.Value != "new" {
		t.Fatalf("WithValue value = %q, want %q", got.Value, "new")
	}
}

func TestStampedSnapshot(t *testing.T) {
	stamped := Stamped[string]{Revision: Revision(4), Updated: time.Unix(10, 0), Value: "value"}
	snap := stamped.Snapshot()

	if snap.Revision != stamped.Revision {
		t.Fatalf("Snapshot revision = %d, want %d", snap.Revision, stamped.Revision)
	}
	if snap.Value != stamped.Value {
		t.Fatalf("Snapshot value = %q, want %q", snap.Value, stamped.Value)
	}
}
