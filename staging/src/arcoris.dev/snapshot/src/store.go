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
	"sync"
	"time"

	"arcoris.dev/chrono/clock"
)

// Store is a concurrency-safe holder for one always-present mutable value.
//
// Store owns its internal value and uses CloneFunc to isolate writes and reads.
// It is the safe baseline for state that may contain slices, maps, pointers, or
// other mutable data. Snapshot and Stamped return cloned values; mutating a value
// returned from Store must not affect Store's internal state when the CloneFunc is
// correct.
//
// Store starts at revision 1 because NewStore commits the initial value. Use a
// value-level container such as maybe.Maybe[T] when the logical state can be
// absent.
//
// Store must be created with NewStore. The zero value is invalid because Store
// requires a CloneFunc and an initial committed value.
//
// Store is safe for concurrent use. Store must not be copied after first use.
type Store[T any] struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects value, revision, and updated.
	mu sync.RWMutex

	// clone copies values across the Store ownership boundary.
	clone CloneFunc[T]

	// clock provides local commit timestamps for Stamped values.
	clock clock.PassiveClock

	// value is the currently committed internal value owned by the Store.
	value T

	// revision is the source-local revision of value.
	revision Revision

	// updated is the local time at which value was committed.
	updated time.Time
}
