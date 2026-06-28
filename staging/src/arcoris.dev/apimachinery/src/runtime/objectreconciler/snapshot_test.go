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

package objectreconciler

import "testing"

func TestSnapshotFromSourcePreservesRevisionAndView(t *testing.T) {
	source := readySourceSnapshot(t, 12)

	converted := snapshotFromSource(source)
	if converted.Revision != source.Revision {
		t.Fatalf("Revision = %s; want %s", converted.Revision, source.Revision)
	}
	if converted.View.Revision() != source.Value.Revision() {
		t.Fatalf("View.Revision() = %s; want %s", converted.View.Revision(), source.Value.Revision())
	}
	if converted.View.Len() != source.Value.Len() {
		t.Fatalf("View.Len() = %d; want %d", converted.View.Len(), source.Value.Len())
	}
}

func TestSnapshotFromSourcePreservesRevisionInvariant(t *testing.T) {
	converted := snapshotFromSource(readySourceSnapshot(t, 7))

	if converted.Revision != converted.View.Revision() {
		t.Fatalf("Revision = %s, View.Revision() = %s", converted.Revision, converted.View.Revision())
	}
}
