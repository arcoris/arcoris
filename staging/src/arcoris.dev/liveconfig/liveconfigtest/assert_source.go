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

// RequireSourceRevision fails the test when src reports a different revision
// than want.
func RequireSourceRevision(t TestingT, src snapshot.RevisionSource, want snapshot.Revision) {
	t.Helper()

	if got := src.Revision(); got != want {
		t.Fatalf("source revision = %d, want %d", got, want)
	}
}

// RequireSourceValue fails the test when src.Snapshot().Value differs from want.
func RequireSourceValue[T any](t TestingT, src snapshot.Source[T], want T, equal func(T, T) bool) {
	t.Helper()
	RequireSnapshotValue(t, src.Snapshot(), want, equal)
}

// RequireConfigSourceValue fails the test when src.Snapshot().Value differs from
// want according to EqualConfig.
func RequireConfigSourceValue(t TestingT, src snapshot.Source[Config], want Config) {
	t.Helper()
	RequireSourceValue(t, src, want, EqualConfig)
}

// RequireStampedNonZeroRevision fails the test when stamped has ZeroRevision.
func RequireStampedNonZeroRevision[T any](t TestingT, stamped snapshot.Stamped[T]) {
	t.Helper()

	if stamped.IsZeroRevision() {
		t.Fatalf("stamped revision is zero")
	}
}

// RequireStampedValue fails the test when stamped.Value differs from want.
func RequireStampedValue[T any](t TestingT, stamped snapshot.Stamped[T], want T, equal func(T, T) bool) {
	t.Helper()
	RequireValue(t, stamped.Value, want, equal)
}

// RequireConfigStampedValue fails the test when stamped.Value differs from want
// according to EqualConfig.
func RequireConfigStampedValue(t TestingT, stamped snapshot.Stamped[Config], want Config) {
	t.Helper()
	RequireStampedValue(t, stamped, want, EqualConfig)
}
