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

import (
	"testing"

	panicassert "arcoris.dev/testutil/panic"
)

func TestMapPreservesRevision(t *testing.T) {
	snap := Snapshot[testRevision, string]{
		Revision: testRevision("domain-r3"),
		Value:    "value",
	}

	got := Map(snap, func(v string) int {
		return len(v)
	})

	if got.Revision != snap.Revision {
		t.Fatalf("Map revision = %q, want %q", got.Revision, snap.Revision)
	}
	if got.Value != 5 {
		t.Fatalf("Map value = %d, want 5", got.Value)
	}
}

func TestMapPanicsOnNilFunction(t *testing.T) {
	snap := Snapshot[testRevision, string]{Revision: testRevision("domain-r4"), Value: "value"}

	panicassert.RequireMessage(t, "snapshot: nil map function", func() {
		_ = Map[testRevision, string, int](snap, nil)
	})
}
