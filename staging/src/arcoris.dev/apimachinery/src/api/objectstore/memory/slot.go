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

package memory

import "sync/atomic"

// slot owns the current immutable record pointer for one object key.
//
// It intentionally uses sync/atomic.Pointer rather than padded atomic
// primitives because there can be many object slots. Compact per-object state
// is preferable until benchmarks prove false sharing dominates memory cost.
type slot struct {
	// current points to the latest immutable live or tombstone record.
	current atomic.Pointer[record]
}

// load returns the currently published record pointer.
func (s *slot) load() *record {
	return s.current.Load()
}

// compareAndSwap publishes next only if current is still old.
func (s *slot) compareAndSwap(old *record, next *record) bool {
	return s.current.CompareAndSwap(old, next)
}
