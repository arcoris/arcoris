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

package atomicx

// Uint64Counter is a padded monotonic uint64 lifetime counter.
//
// Uint64Counter represents event-like runtime accounting that only moves
// forward. Typical examples include:
//
//   - admitted requests;
//   - rejected requests;
//   - completed work items;
//   - failed work items;
//   - retry attempts;
//   - dropped events;
//   - cache hits;
//   - cache misses;
//   - controller ticks;
//   - dispatch attempts;
//   - successful state transitions;
//   - failed state transitions.
//
// Uint64Counter is the default counter type for ARCORIS component accounting.
// Prefer it over Uint32Counter unless a 32-bit counter is part of an explicit
// bounded state model or external protocol boundary.
//
// Uint64Counter is not a gauge. It intentionally does not expose Store, Swap,
// Sub, Dec, or CompareAndSwap methods. A lifetime counter should not be reset,
// decremented, or conditionally rewritten by ordinary runtime code. Activity
// windows, recent views, and reporting intervals should be built from separate
// snapshot/delta value types instead of mutating the source counter.
//
// Uint64Counter uses PaddedUint64 internally so it can be embedded in hot
// component/runtime structs with reduced risk of false sharing.
//
// Counter arithmetic follows ordinary uint64 atomic arithmetic. If a long-lived
// counter wraps from the largest uint64 value to zero, unsigned snapshot/delta
// logic can account for one wrap between two samples. Multiple wraps between two
// samples cannot be detected from two uint64 values alone and must be avoided by
// reasonable sampling cadence.
//
// Uint64Counter is zero-value usable.
//
// Uint64Counter must not be copied after first use. Copying a live counter can
// split one logical lifetime counter into independent copies and corrupt runtime
// accounting. Construct it in place, pass it by pointer when sharing, and do not
// copy containing structs after the counter becomes active.
type Uint64Counter struct {
	noCopy noCopy
	value  PaddedUint64
}

// Load atomically returns the current lifetime counter value.
//
// Load observes exactly one atomic value. It does not make a multi-field
// accounting snapshot globally consistent. If a caller needs a consistent view
// of multiple counters or gauges, the caller must provide synchronization at the
// owner level.
//
// Load is useful for diagnostics, direct assertions, and low-level reads. Use
// snapshot/delta value types for windowed activity measurement.
func (c *Uint64Counter) Load() uint64 {
	return c.value.Load()
}

// Add atomically adds delta events to the counter and returns the new lifetime
// value.
//
// Add is intended for batched event increments. For the common single-event
// case, use Inc.
//
// Add follows ordinary uint64 arithmetic and may wrap. Wraparound is allowed for
// lifetime counters because deltas can be computed with unsigned modulo
// arithmetic. If the value represents current state that must never wrap, use
// Uint64Gauge instead.
//
// Add accepts zero. Adding zero leaves the logical counter unchanged and returns
// the currently observed value.
func (c *Uint64Counter) Add(delta uint64) uint64 {
	return c.value.Add(delta)
}

// Inc atomically adds one event to the counter and returns the new lifetime
// value.
//
// Inc is a convenience wrapper around Add(1). Use Inc for the common case where
// one completed event, rejected event, cache hit, controller tick, or similar
// unit of work is recorded.
func (c *Uint64Counter) Inc() uint64 {
	return c.value.Inc()
}
