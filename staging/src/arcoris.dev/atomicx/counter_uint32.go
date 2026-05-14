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

// Uint32Counter is a padded monotonic uint32 lifetime counter.
//
// Uint32Counter represents event-like runtime accounting that only moves
// forward, but only for counters that deliberately use a 32-bit range. Typical
// valid use cases are narrow, bounded, or protocol-shaped counters where uint32
// is part of the intended state model.
//
// For general runtime event accounting, prefer Uint64Counter. A uint32 counter
// wraps much sooner than a uint64 counter, and explicit cache-line padding means
// the memory saving is usually not meaningful for hot shared fields.
//
// Uint32Counter may be appropriate for:
//
//   - compact bounded component-local event counters;
//   - hot counters that mirror an external uint32 protocol field;
//   - short-lived counters with controlled sampling cadence;
//   - test or simulation counters where a smaller wrap boundary is intentional.
//
// Uint32Counter is not a gauge. It intentionally does not expose Store, Swap,
// Sub, or Dec methods. A lifetime counter should not be reset or decremented by
// ordinary runtime code. Workload windows, recent activity views, and reporting
// intervals should compute deltas between snapshots instead of mutating the
// source counter.
//
// Uint32Counter uses PaddedUint32 internally so it can be embedded in hot
// component/runtime structs with reduced risk of false sharing.
//
// Counter arithmetic follows ordinary uint32 atomic arithmetic. If the counter
// wraps from the largest uint32 value to zero, Uint32CounterDelta remains correct
// for a single wrap between two samples. Multiple wraps between two samples
// cannot be detected from two uint32 values alone and must be avoided by
// sampling frequently enough for the expected event rate.
//
// Uint32Counter is zero-value usable.
//
// Uint32Counter must not be copied after first use. Copying a live counter can
// split one logical lifetime counter into independent copies and corrupt runtime
// accounting. Construct it in place, pass it by pointer when sharing, and do not
// copy containing structs after the counter becomes active.
type Uint32Counter struct {
	noCopy noCopy
	value  PaddedUint32
}

// Load atomically returns the current lifetime counter value.
//
// Load observes exactly one atomic value. It does not make a multi-field
// accounting snapshot globally consistent. If a caller needs a consistent view
// of multiple counters or gauges, the caller must provide synchronization at the
// owner level.
func (c *Uint32Counter) Load() uint32 {
	return c.value.Load()
}

// Add atomically adds delta events to the counter and returns the new lifetime
// value.
//
// Add is intended for batched event increments. For the common single-event
// case, use Inc.
//
// Add follows ordinary uint32 arithmetic and may wrap. Wraparound is allowed for
// lifetime counters because deltas are computed with unsigned modulo arithmetic.
// If the value represents current state that must never wrap, use Uint32Gauge or
// Uint64Gauge instead.
//
// Prefer Uint64Counter unless the 32-bit range is an explicit part of the
// counter's design.
func (c *Uint32Counter) Add(delta uint32) uint32 {
	return c.value.Add(delta)
}

// Inc atomically adds one event to the counter and returns the new lifetime
// value.
//
// Inc is a convenience wrapper around Add(1).
func (c *Uint32Counter) Inc() uint32 {
	return c.value.Inc()
}
