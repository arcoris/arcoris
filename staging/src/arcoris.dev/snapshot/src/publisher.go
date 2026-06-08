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

	"arcoris.dev/atomicx"
	"arcoris.dev/chrono/clock"
)

// Publisher atomically publishes immutable copy-on-write values.
//
// Publisher is the fast baseline for read-mostly values such as observer lists,
// handler registries, immutable routing tables, and precomputed dispatch plans.
// Reads use an atomic pointer load and do not clone the value.
//
// Writes are serialized so each commit advances the source-local revision and
// stores the corresponding record as one ordered publication. Reads remain
// lock-free by loading the latest published record pointer.
//
// Publisher does not freeze, clone, or deep-copy values. Values passed to
// Publish or PublishStamped must be treated as immutable after publication. If T
// contains slices, maps, pointers, or other mutable state, the caller must build
// a fresh copy before publishing it and must not mutate that copy after Publish
// returns. Use Store when caller-owned clone isolation is required.
//
// Publisher is safe for concurrent use and zero-value usable. A zero-value
// Publisher has no published record and returns a zero snapshot until the first
// Publish. Publisher must not be copied after first use.
type Publisher[T any] struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu serializes Publish and PublishStamped commits.
	//
	// The mutex preserves monotonic visible revisions under concurrent
	// publishers. Without it, one goroutine could reserve a newer revision and
	// store it before another goroutine stores an older reserved revision.
	mu sync.Mutex

	// nextRevision is the latest source-local revision assigned to a published
	// record.
	//
	// nextRevision is protected by mu. Readers do not read this field; they load
	// the revision from the latest immutable record through ptr.
	nextRevision Revision

	// clock provides local publication timestamps for Stamped values.
	//
	// A nil clock means the zero-value Publisher should lazily use RealClock.
	clock clock.PassiveClock

	// ptr is a padded atomic pointer to the latest immutable published record.
	//
	// The padding isolates the hot pointer cell from neighboring fields. It does
	// not provide publication ordering; write-side ordering is provided by mu.
	ptr atomicx.PaddedPointer[record[T]]
}

// record is an immutable published value.
//
// Once a record is stored in Publisher.ptr, it must never be mutated. Readers may
// load and use records without locks.
type record[T any] struct {
	// revision is the source-local revision assigned at publication time.
	revision Revision

	// updated is the local publication time.
	updated time.Time

	// value is the immutable published value.
	value T
}
