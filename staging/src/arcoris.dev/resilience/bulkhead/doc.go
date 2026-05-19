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

// Package bulkhead provides local concurrency isolation for ARCORIS component
// internals.
//
// A bulkhead bounds the number of operations that may execute inside one
// protected section at the same time. Callers acquire a Permit before entering
// the protected section and release the Permit when the section is left. When no
// capacity is available, the limiter rejects the admission attempt immediately.
//
// The package models a local, non-blocking bulkhead. It does not queue waiters,
// park goroutines, schedule work, retry operations, classify operation errors,
// enforce request deadlines, rate-limit throughput, perform circuit breaking,
// coordinate distributed capacity, export metrics, or integrate with health.
// Those concerns belong to packages layered above this primitive.
//
// The first implementation intentionally exposes TryAcquire only. This keeps
// the package focused on concurrency isolation and avoids embedding fairness,
// waiter lifecycle, cancellation, wake-up ordering, or queue scheduling policy in
// the base bulkhead. Future packages may build blocking, queued, keyed, or
// dynamically reconfigured bulkheads on top of this foundation.
//
// A Limiter owns coherent mutable state under an internal mutex and publishes a
// revisioned read model through package snapshot. Snapshot reads are delegated to
// snapshot.Publisher and do not acquire the limiter mutex. Permits are
// release-once capabilities: releasing the same Permit more than once is a
// no-op and cannot underflow the limiter.
//
// Revisions are source-local publication versions inherited from package
// snapshot. They are useful for cheap change detection by consumers of the same
// limiter, but they are not a global ordering across different limiters or
// components.
package bulkhead
