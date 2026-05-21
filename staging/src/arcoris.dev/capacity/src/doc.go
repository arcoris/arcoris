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

// Package capacity provides local scalar capacity accounting primitives for
// ARCORIS component internals.
//
// The package owns capacity limits, reserved capacity, available capacity,
// capacity debt after limit reduction, reservation ownership, and consistent
// revisioned snapshots of local capacity state. It is intentionally below
// admission, scheduling, rate limiting, overload control, and worker-isolation
// policy.
//
// # Model
//
// A Ledger owns one local capacity limit and the amount currently reserved from
// that limit. TryReserve is a non-blocking check-and-reserve operation. When it
// succeeds, it returns a Reservation that owns the reserved amount until Release
// or TryRelease returns that amount to the ledger.
//
// Limit changes never revoke existing reservations. If a limit is reduced below
// already reserved capacity, the ledger reports capacity debt and refuses new
// reservations until releases bring reserved capacity back under the current
// limit.
//
// # Boundaries
//
// Capacity does not provide blocking Acquire operations, wait queues, context
// cancellation, fairness policy, rate limiting, adaptive concurrency control,
// admission decisions, worker pools, bulkheads, circuit breakers, health gates,
// logging, metrics, tracing, distributed coordination, or multi-resource
// scheduling.
//
// Higher-level packages may build those behaviors on top of Ledger,
// Reservation, Amount, and Snapshot, but this package must remain a small local
// accounting layer.
//
// # Dependency policy
//
// Production code in this package depends only on the Go standard library and
// arcoris.dev/snapshot for source-local revisioned read models.
package capacity
