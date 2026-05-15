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

// Package atomicx provides padded atomic accounting primitives for ARCORIS
// component internals.
//
// The package is intended for hot runtime state that is shared across goroutines
// and updated frequently enough for false sharing and accounting invariants to
// matter. Typical users are schedulers, admission controllers, queues, worker
// runtimes, caches, dispatch loops, adaptive controllers, and component-local
// state machines.
//
// This package is the root of arcoris.dev/atomicx, not apimachinery. It does not
// define API object contracts, schema identities, metadata, status objects, resource
// versions, or external observability formats. It provides low-level concurrent
// state containers and copyable sampling values that higher-level packages may
// use to build component accounting, metrics, diagnostics, and control loops.
//
// # Scope
//
// atomicx contains four groups of types:
//
//   - raw padded atomics;
//   - monotonic unsigned counters;
//   - counter snapshots and deltas;
//   - current-state gauges.
//
// Raw padded atomics are the lowest-level primitives:
//
//   - PaddedUint64;
//   - PaddedUint32;
//   - PaddedInt64;
//   - PaddedInt32.
//
// They expose raw atomic operations and intentionally do not enforce semantic
// accounting rules. For example, PaddedUint64.Add uses ordinary unsigned
// arithmetic and may wrap. PaddedInt64.Add follows atomic.Int64 semantics and
// does not reject signed overflow or underflow. Raw padded atomics are storage
// cells, not domain-aware accounting types.
//
// Counters are monotonic lifetime event counters:
//
//   - Uint64Counter;
//   - Uint32Counter.
//
// They are intended for values that only move forward, such as admitted
// requests, rejected requests, completed work items, failed work items, cache
// hits, cache misses, dropped events, retry attempts, dispatch attempts, and
// controller ticks.
//
// Counters intentionally do not expose Store, Swap, Sub, Dec, or
// CompareAndSwap. A lifetime counter should not be reset, decremented, or
// conditionally rewritten by ordinary runtime code. Recent activity windows
// should be computed from snapshots and deltas instead of mutating the source
// counter.
//
// Counter snapshots and deltas are copyable value objects:
//
//   - Uint64CounterSnapshot;
//   - Uint32CounterSnapshot;
//   - Uint64CounterDelta;
//   - Uint32CounterDelta.
//
// Snapshots capture one counter value at one point in time. Deltas describe the
// unsigned difference between two counter samples. Delta calculation is
// wrap-aware for one unsigned wrap between two samples. Multiple wraps between
// two samples cannot be detected from two values alone and must be prevented by
// choosing a reasonable sampling cadence.
//
// Gauges represent current runtime state that can move both up and down:
//
//   - Uint64Gauge;
//   - Uint32Gauge;
//   - Int64Gauge;
//   - Int32Gauge.
//
// Unsigned gauges are for non-negative current quantities such as queue depth,
// in-flight operation count, retained bytes, active leases, pending work,
// reserved capacity, and available permits. They reject unsigned overflow and
// underflow.
//
// Signed gauges are for current quantities where negative values are meaningful,
// such as correction deltas, signed budget movement, controller drift, signed
// reconciliation offsets, capacity debt, and temporary signed runtime
// adjustments. They reject signed overflow and underflow.
//
// # File ownership
//
// atomicx intentionally separates files by responsibility and numeric type.
//
// Raw padded primitives live in:
//
//   - uint64.go;
//   - uint32.go;
//   - int64.go;
//   - int32.go.
//
// Mutable counters live in:
//
//   - counter_uint64.go;
//   - counter_uint32.go.
//
// Counter sampling values live in:
//
//   - snapshot_uint64.go;
//   - snapshot_uint32.go;
//   - delta_uint64.go;
//   - delta_uint32.go.
//
// Gauges live in:
//
//   - gauge_uint64.go;
//   - gauge_uint32.go;
//   - gauge_int64.go;
//   - gauge_int32.go.
//
// Counter files must not define snapshots or deltas. Snapshot files must not
// define mutable counters. Delta files must not define counter state. Gauge
// files must keep their panic constants and numeric bounds local to the gauge
// implementation that uses them.
//
// # Non-goals
//
// atomicx is not a metrics package. It does not define metric names, labels,
// descriptors, registries, exporters, histograms, summaries, Prometheus
// collectors, OpenTelemetry instruments, or scrape endpoints.
//
// atomicx is not an API machinery package. It must not contain schema types,
// TypeMeta, ObjectMeta, runtime objects, serializers, status objects, API
// errors, selectors, resource versions, or Kubernetes compatibility adapters.
//
// atomicx is not a domain package. It must not contain scheduler policies,
// admission algorithms, queue implementations, worker protocols, bufferpool
// classes, shard counters, workload classes, tenant IDs, or ARCORIS control-plane
// resource definitions.
//
// atomicx does not mirror every type in sync/atomic. New padded atomic types
// should be added only when ARCORIS components have a concrete hot-path need for
// that representation and the package can document the copy, overflow,
// ownership, and synchronization rules clearly.
//
// # Counter versus gauge semantics
//
// Counters and gauges are intentionally separate.
//
// A counter is a lifetime event count. It moves forward, may wrap according to
// unsigned arithmetic, and is sampled through snapshots and deltas. Counters are
// appropriate for historical activity:
//
//	var admitted atomicx.Uint64Counter
//	admitted.Inc()
//
// A gauge is current state. It can move up and down and enforces bounds.
// Gauges are appropriate for live quantities:
//
//	var inFlight atomicx.Uint64Gauge
//	inFlight.Inc()
//	defer inFlight.Dec()
//
// Signed values are gauges, not counters. If a value can become negative, it is
// current signed state, correction, drift, debt, or movement. It should use a
// signed gauge or a raw signed padded atomic, not a signed counter.
//
// # Raw padded atomics
//
// Raw padded atomics are deliberately low-level. They expose atomic operations
// without domain-level invariants. They are useful when a caller owns a custom
// state model, such as a compact state word, bitmask, explicit state-machine
// transition, or owner-controlled publication protocol.
//
// Prefer semantic wrappers when the value has common accounting meaning:
//
//   - use Uint64Counter for general lifetime event accounting;
//   - use Uint32Counter only when a 32-bit counter is deliberate;
//   - use Uint64Gauge for non-negative current quantities;
//   - use Uint32Gauge only when a bounded 32-bit gauge is deliberate;
//   - use Int64Gauge for general signed current quantities;
//   - use Int32Gauge only when a bounded signed 32-bit gauge is deliberate.
//
// # Padding and false sharing
//
// Padded primitives use an explicit leading pad before the atomic value and a
// trailing line-completion pad after the value. The leading pad separates the
// value from preceding fields. The trailing completion fills the rest of the
// value's 64-byte slot so following fields and slice elements do not share that
// slot.
//
// CacheLinePadSize is an ARCORIS layout policy, not a runtime-detected hardware
// cache-line size.
//
// False sharing can occur when unrelated variables occupy the same CPU cache
// line and are updated by different goroutines on different cores. Even when the
// logical variables are independent, writes to one variable can invalidate the
// cache line containing another variable. This can create unnecessary coherence
// traffic and degrade hot-path performance.
//
// Padding is not free. It deliberately increases memory footprint. Padded
// atomics should be used for a small number of hot shared fields, not for
// ordinary per-object metadata, API objects, per-request structs, per-item
// storage, or cold state.
//
// # Zero values
//
// Mutable atomicx state containers are zero-value usable. A zero-value counter,
// gauge, or raw padded atomic is ready to use without explicit initialization.
//
// This property is important because component-local accounting structs should
// be embeddable directly into larger runtime structs without constructor-heavy
// setup code.
//
// # Copying
//
// Mutable atomicx state containers must not be copied after first use.
//
// This rule applies to:
//
//   - PaddedUint64;
//   - PaddedUint32;
//   - PaddedInt64;
//   - PaddedInt32;
//   - Uint64Counter;
//   - Uint32Counter;
//   - Uint64Gauge;
//   - Uint32Gauge;
//   - Int64Gauge;
//   - Int32Gauge.
//
// Copying a live atomic state container splits one logical state cell into
// independent copies. That can make events disappear, duplicate current-state
// accounting, corrupt capacity tracking, or invalidate ownership assumptions.
//
// Mutable state containers contain an internal noCopy marker so tools such as
// go vet -copylocks can detect accidental copies. The marker is a static
// analysis aid only. It does not provide runtime protection.
//
// Snapshot and delta types are intentionally copyable. They are immutable value
// objects and must not contain noCopy.
//
// # Snapshot consistency
//
// Atomic Load observes one atomic value. It does not make a larger multi-field
// struct globally consistent.
//
// A collection of counter snapshots taken from several counters is not an
// atomic snapshot of the entire component unless the caller provides additional
// synchronization around the sampling operation.
//
// atomicx provides per-value atomicity. Component-level consistency belongs to
// the owner of the surrounding runtime state.
//
// # Panic and Try semantics
//
// Gauge Add, Sub, Inc, and Dec methods enforce accounting invariants and panic
// when an operation would cross the valid numeric boundary.
//
// For unsigned gauges:
//
//   - Add and Inc panic on overflow;
//   - Sub and Dec panic on underflow.
//
// For signed gauges:
//
//   - Add and Sub panic on signed overflow or underflow;
//   - Inc and Dec inherit the corresponding Add/Sub boundary checks.
//
// These panics indicate internal accounting bugs. They are appropriate when a
// component must never release more state than it acquired, reserve more state
// than the numeric range can represent, or silently wrap a current-state value.
//
// TryAdd and TrySub are provided for control paths where refusal is expected
// behavior. Examples include capacity reservation, best-effort permit
// acquisition, admission checks, and bounded budget movement. On failure, TryAdd
// and TrySub leave the gauge unchanged and return the current value observed at
// the failing attempt together with false.
//
// # Serialization
//
// Mutable atomicx types are runtime state containers and should not be serialized
// directly. They are not API objects and are not stable wire-format values.
//
// If external reporting is needed, export loaded values, snapshots, deltas, or
// higher-level metric samples produced by a metrics package. Do not marshal live
// counters or gauges as API state.
//
// # Stability
//
// atomicx is a focused ARCORIS module. Public APIs in this package should be
// treated as deliberately designed, but ARCORIS may still refine package
// boundaries before the first stable release.
//
// # Dependency policy
//
// atomicx should remain dependency-light. Production code in this package should
// depend only on the Go standard library.
//
// The package must not depend on apimachinery, scheduler packages, queue
// packages, component implementations, observability exporters, or external
// assertion libraries.
package atomicx
