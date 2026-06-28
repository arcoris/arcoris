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

func TestStoreImplementsSourceInterfaces(t *testing.T) {
	var _ Source[LocalRevision, string] = (*Store[string])(nil)
	var _ StampedSource[LocalRevision, string] = (*Store[string])(nil)
	var _ RevisionSource[LocalRevision] = (*Store[string])(nil)
}

func TestPublisherImplementsSourceInterfaces(t *testing.T) {
	var _ Source[LocalRevision, string] = (*Publisher[string])(nil)
	var _ StampedSource[LocalRevision, string] = (*Publisher[string])(nil)
	var _ RevisionSource[LocalRevision] = (*Publisher[string])(nil)
}

func TestSourceReadsStore(t *testing.T) {
	store := NewStore("value", IdentityClone[string])
	read := func(src Source[LocalRevision, string]) Snapshot[LocalRevision, string] {
		return src.Snapshot()
	}

	if got, want := read(store).Value, "value"; got != want {
		t.Fatalf("source value = %q, want %q", got, want)
	}
}
