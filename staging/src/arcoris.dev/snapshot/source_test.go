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

import "testing"

func TestStoreImplementsSourceInterfaces(t *testing.T) {
	var _ Source[string] = (*Store[string])(nil)
	var _ StampedSource[string] = (*Store[string])(nil)
	var _ RevisionSource = (*Store[string])(nil)
}

func TestPublisherImplementsSourceInterfaces(t *testing.T) {
	var _ Source[string] = (*Publisher[string])(nil)
	var _ StampedSource[string] = (*Publisher[string])(nil)
	var _ RevisionSource = (*Publisher[string])(nil)
}

func TestSourceReadsStore(t *testing.T) {
	store := NewStore("value", Identity[string])
	read := func(src Source[string]) Snapshot[string] {
		return src.Snapshot()
	}

	if got, want := read(store).Value, "value"; got != want {
		t.Fatalf("source value = %q, want %q", got, want)
	}
}
