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

package liveconfig

import "arcoris.dev/snapshot"

// Change describes the result of applying a candidate configuration.
//
// Previous is the snapshot that was current before Apply evaluated the
// candidate. Current is the snapshot that remains current after Apply returns.
// When Changed is false, Current and Previous refer to the same source revision.
// When Changed is true, Current is the newly published snapshot.
//
// A rejected candidate returns Changed=false with Current equal to Previous and
// a non-nil error from Apply. A valid no-op also returns Changed=false, but with
// a nil error because the candidate was accepted and intentionally not
// published.
type Change[T any] struct {
	// Previous is the snapshot visible before the candidate was evaluated.
	Previous snapshot.Snapshot[T]

	// Current is the snapshot visible after the candidate was evaluated.
	Current snapshot.Snapshot[T]

	// Changed reports whether Apply published a new source revision.
	Changed bool
}

// IsChanged reports whether Apply published a new source revision.
func (c Change[T]) IsChanged() bool {
	return c.Changed
}

// IsNoop reports whether Apply left the current source revision unchanged.
func (c Change[T]) IsNoop() bool {
	return !c.Changed
}
