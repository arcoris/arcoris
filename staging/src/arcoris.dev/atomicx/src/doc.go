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

// Package atomicx provides padded atomic primitives for hot ARCORIS runtime
// accounting fields.
//
// The package is intentionally small. It is not a generic mirror of sync/atomic
// and does not provide every integer width. Numeric atomicx APIs are 64-bit
// only because padded values are meant for a small number of independently hot
// runtime fields, and the fixed padding dominates the memory layout. Smaller
// integer widths do not provide meaningful memory savings in this padded memory
// model. Compact per-object atomic state is a different problem; callers should
// use sync/atomic directly or a future non-padded package for that use case.
//
// atomicx contains four public categories:
//
//   - raw padded cells: PaddedUint64, PaddedInt64, and PaddedPointer;
//   - monotonic lifetime counters: Uint64Counter;
//   - copyable counter samples and deltas: Uint64CounterSample and
//     Uint64CounterDelta;
//   - checked current-state gauges: Uint64Gauge and Int64Gauge.
//
// Raw padded numeric cells expose low-level atomic operations such as Load,
// Store, Add, Swap, and CompareAndSwap. They do not enforce accounting
// invariants. Uint64Counter narrows that raw layer into monotonic event
// accounting and intentionally omits reset, decrement, and CAS-style APIs.
// Gauges represent current state and enforce overflow or underflow checks during
// arithmetic.
//
// PaddedPointer remains in the package because pointer publication is a
// separate atomic category from numeric accounting. It is a raw pointer cell: it
// does not own, clone, freeze, retain, release, version, reclaim, or protect the
// pointed object. Callers own immutability, lifetime, ABA, and reclamation
// rules.
//
// Mutable atomicx containers are zero-value usable and must not be copied after
// first use. Counter samples and deltas are immutable value objects and are safe
// to copy. Atomic operations provide per-value atomicity only; callers that need
// a coherent view across multiple fields must provide owner-level
// synchronization.
//
// Production code in this module depends only on the Go standard library.
package atomicx
