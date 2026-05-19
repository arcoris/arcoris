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

// Package bulkhead provides bounded in-flight isolation for ARCORIS component
// internals.
//
// A Bulkhead is a small resilience-domain wrapper around capacity.Ledger. It
// reserves one local capacity unit before protected work starts and returns that
// unit when the Lease is released. When no capacity is available, TryAcquire
// rejects immediately.
//
// The package is local and non-blocking. It does not wait, queue callers,
// observe contexts, implement fairness, rate-limit throughput, retry work,
// classify operation failures, integrate with health, export metrics, log, trace,
// schedule work, or manage worker pools. Those policies belong above this
// primitive.
//
// capacity owns scalar accounting and lease/release semantics. bulkhead owns the
// resilience meaning: bounded in-flight isolation around protected execution.
// The public read model is intentionally the capacity snapshot itself so callers
// see the same Limit, Reserved, Available, and Debt semantics at both layers.
package bulkhead
