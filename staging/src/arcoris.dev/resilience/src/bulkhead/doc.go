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

// Package bulkhead provides bounded in-flight isolation for ARCORIS component
// internals.
//
// A Bulkhead is a small resilience-domain wrapper around capacity.Ledger. It
// reserves local capacity before protected work starts and returns that capacity
// when the Lease is released. When no capacity is available, acquisition rejects
// immediately.
//
// TryAcquire and TryAcquireAmount are the direct APIs. A successful acquisition
// returns a Lease and an Observation with RefusalNone. A denied acquisition
// returns no lease and an Observation that preserves both the observed capacity
// snapshot and the precise capacity refusal. Successful direct acquisition
// allocates only the bulkhead Lease token; it does not also allocate a
// capacity.Reservation object.
//
// The package is local and non-blocking. It does not wait, queue callers,
// observe contexts, implement fairness, rate-limit throughput, retry work,
// classify operation failures, integrate with health, export metrics, log, trace,
// schedule work, map admission reasons, use admission metadata catalogs, or
// manage worker pools. Those policies belong above this primitive.
//
// capacity owns scalar accounting, limits, debt, revisions, and refusal
// taxonomy. bulkhead owns the resilience meaning: protected in-flight work,
// release-once Lease ownership, and direct Observation metadata. Snapshots are
// capacity-only and intentionally do not include queues, waiters, fairness,
// metrics, or historical counters.
//
// SetLimit never revokes active leases. Lowering the limit below active leases
// creates debt in the underlying capacity snapshot; new acquisitions are denied
// until releases or a later limit increase restore availability.
//
// Package arcoris.dev/resilience/bulkheadadmission maps direct bulkhead
// observations into admission.Result values for callers that need the generic
// admission surface.
package bulkhead
