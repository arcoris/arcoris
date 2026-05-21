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

package liveconfigtest

import "arcoris.dev/snapshot"

// RequireNonZeroRevision fails the test when snap has ZeroRevision.
func RequireNonZeroRevision[T any](t TestingT, snap snapshot.Snapshot[T]) {
	t.Helper()

	if snap.IsZeroRevision() {
		t.Fatalf("snapshot revision is zero")
	}
}

// RequireRevision fails the test when snap has a different revision than want.
func RequireRevision[T any](t TestingT, snap snapshot.Snapshot[T], want snapshot.Revision) {
	t.Helper()

	if snap.Revision != want {
		t.Fatalf("snapshot revision = %d, want %d", snap.Revision, want)
	}
}

// RequireChangedSince fails the test when snap has not changed since prev.
func RequireChangedSince[T any](t TestingT, snap snapshot.Snapshot[T], prev snapshot.Revision) {
	t.Helper()

	if !snap.ChangedSince(prev) {
		t.Fatalf("snapshot revision = %d, want changed since %d", snap.Revision, prev)
	}
}

// RequireUnchangedSince fails the test when snap has changed since prev.
func RequireUnchangedSince[T any](t TestingT, snap snapshot.Snapshot[T], prev snapshot.Revision) {
	t.Helper()

	if snap.ChangedSince(prev) {
		t.Fatalf("snapshot revision = %d, want unchanged since %d", snap.Revision, prev)
	}
}

// RequireSnapshotValue fails the test when snap.Value differs from want.
func RequireSnapshotValue[T any](t TestingT, snap snapshot.Snapshot[T], want T, equal func(T, T) bool) {
	t.Helper()
	RequireValue(t, snap.Value, want, equal)
}
