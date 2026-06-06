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

// Package capacity provides local resource-accounting foundations for ARCORIS.
//
// The package owns accounting truth: exact amounts, stable resource identities,
// canonical resource vectors, non-empty demands, limits, reserved allocations,
// available capacity, per-resource debt after limit shrink, all-or-nothing
// vector reservation, scalar raw accounting, reservation ownership, and
// revisioned observations.
//
// Ledger is the optimized scalar hot-path owner for common single-resource
// limits such as bulkhead slots, worker slots, queue slots, and active request
// counts. Its raw TryReserve and Release methods do not allocate or build
// snapshots. TryAcquire returns a capacity-owned Reservation when callers want
// ownership handled by this package. Observed methods and Snapshot build
// diagnostics only when explicitly requested. Scalar snapshots are safe
// concurrent observations, not global serialization points for raw accounting.
//
// VectorLedger is the explicit strict owner for multi-resource all-or-nothing
// accounting. Value types are copy-safe and immutable through their public APIs;
// stateful ledgers and reservations are constructor-created owners and must not
// be copied after first use.
//
// Capacity deliberately does not own catalogs, descriptors, admission decisions,
// queues, schedulers, fairness policy, quotas, metrics, health, runtime
// execution, or distributed coordination. Higher layers interpret the local
// accounting facts returned by this package.
package capacity
