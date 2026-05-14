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

// Package snapshot provides typed, revisioned publication primitives for
// component read models.
//
// A snapshot is a point-in-time read view of component state. Components use this
// package when they own mutable internal state but need to expose a stable value
// to readers, diagnostics, health probes, observers, or tests.
//
// Store is the safe baseline for mutable values. It owns one always-present value,
// protects it with a mutex, isolates readers and writers with an explicit
// CloneFunc, and advances a local Revision when the value changes.
//
// Publisher is the fast baseline for immutable copy-on-write values. It publishes
// immutable records through an atomic pointer so readers can load the latest
// snapshot without locking or cloning. Values passed to Publisher must not be
// mutated after publication.
//
// Snapshot is intentionally lightweight and contains only a Revision and a value.
// Stamped adds the local update timestamp for components that need publication
// metadata. Timestamps come from arcoris.dev/chrono/clock.PassiveClock; snapshot
// does not define its own clock abstraction.
//
// Revisions are monotonic and source-local. ZeroRevision means "no committed
// publication", so Store and Publisher panic instead of wrapping revision
// counters back to zero.
//
// This package does not implement cache policy, TTL, staleness, background
// refresh, loading, watch subscriptions, history, diffing, serialization,
// persistence, event sourcing, or generic deep copy. Optional state should be
// modeled as the value itself, for example Store[maybe.Maybe[T]], rather than by
// adding a separate cache holder to this package.
//
// Revision values are local to one source. They are suitable for detecting
// changes from the same Store or Publisher, but they must not be treated as a
// process-wide or distributed ordering across independent sources.
package snapshot
