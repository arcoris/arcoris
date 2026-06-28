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
	"time"

	panicassert "arcoris.dev/testutil/panic"
)

func TestMapStampedPreservesMetadata(t *testing.T) {
	stamped := Stamped[testRevision, string]{
		Revision: testRevision("domain-r3"),
		Updated:  time.Unix(30, 0),
		Value:    "value",
	}

	got := MapStamped(stamped, func(v string) int {
		return len(v)
	})

	if got.Revision != stamped.Revision {
		t.Fatalf("MapStamped revision = %q, want %q", got.Revision, stamped.Revision)
	}
	if !got.Updated.Equal(stamped.Updated) {
		t.Fatalf("MapStamped updated = %s, want %s", got.Updated, stamped.Updated)
	}
	if got.Value != 5 {
		t.Fatalf("MapStamped value = %d, want 5", got.Value)
	}
}

func TestMapStampedPanicsOnNilFunction(t *testing.T) {
	stamped := Stamped[testRevision, string]{
		Revision: testRevision("domain-r4"),
		Updated:  time.Unix(40, 0),
		Value:    "value",
	}

	panicassert.RequireMessage(t, "snapshot: nil stamped map function", func() {
		_ = MapStamped[testRevision, string, int](stamped, nil)
	})
}
