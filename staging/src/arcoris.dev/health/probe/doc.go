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

// Package probe periodically evaluates explicitly configured health targets and
// exposes the latest cached observations to concurrent readers.
//
// # Package scope
//
// Package probe is an optional runtime cache layer over package health. A Runner
// owns a schedule-driven probe loop, calls an Evaluator source for each
// configured target, stores one latest Snapshot per target, and lets readers
// inspect cached observations without executing checks directly. The default
// schedule is fixed; WithInterval is a convenience wrapper for that common case.
//
// The package intentionally does not define health checks, own health
// registries, duplicate target status aggregation, duplicate reason validation,
// change evaluator execution policy, expose HTTP or gRPC endpoints, publish
// metrics, log diagnostics, attach tracing spans, handle process signals, drive
// lifecycle transitions, implement retries or result-dependent backoff, restart
// components, or make admission, routing, target-specific scheduling, or
// workload-control decisions. Advanced schedules may provide jitter, capped
// sequences, or finite sequences. Finite schedule exhaustion ends Run with an
// error.
//
// # Relationship to health
//
// Package health owns the transport-neutral health contracts: Target, Result,
// Report, Registry, Gate, and TargetPolicy. Package probe depends on a small
// local Evaluator interface instead of reading registries or executing checks
// itself. Package health/eval provides the default synchronous evaluator.
//
// # Transport adapters
//
// Package probe does not expose transport behavior. Package health/http adapts
// package health reports to HTTP endpoints and stays independent from package
// probe by default. Future HTTP, gRPC, CLI, or metrics adapters may choose to
// consume Snapshot values, but those adapters own their wire contracts,
// rendering rules, status mappings, diagnostics, authentication, and exposure
// safety.
//
// # Runtime ownership
//
// NewRunner validates evaluator, clock, schedule, staleness window, and explicit
// target list, but it does not start background work. Run owns exactly one
// schedule sequence while the caller's context is active. Concurrent Run calls
// on the same Runner are rejected. When Run exits after context cancellation,
// the Runner can be started again by the same owner or a later owner.
//
// # Snapshot model
//
// Snapshot is the read model for one target's latest cached observation. The
// zero Snapshot is valid and means no cached observation. Every observed
// Snapshot has a concrete Target, a valid health.Report whose Target matches the
// Snapshot Target, a non-zero Updated timestamp, and a positive per-target
// Revision. Internally, store uses one snapshot.Store[observation] per observed
// target; the underlying snapshot store assigns Revision and Updated, while
// Stale remains read-boundary metadata computed by Runner.
//
// # Staleness
//
// Staleness is cache metadata computed at the read boundary. It is not stored as
// durable cache state. A zero stale-after window disables stale detection.
// Positive windows mark a snapshot stale only when its age is greater than the
// configured window; an age equal to the window is still fresh. Negative ages are
// treated as not stale so mutable test clocks and clock skew do not manufacture
// stale results.
//
// # Concurrency
//
// Runner read methods may be called concurrently with Run. Snapshot and
// Snapshots return detached values. The private store preserves configured
// target order for Snapshots, protects its target map with a mutex, and delegates
// per-target value, revision, timestamp, and clone isolation to snapshot.Store.
//
// # File ownership
//
//   - option.go owns the Option contract and option application order;
//   - config.go owns normalized Runner construction settings;
//   - option_clock.go and option_clock_error.go own clock configuration;
//   - schedule.go and schedule_error.go own probe cadence and initial-probe
//     configuration;
//   - option_stale.go and stale_error.go own stale-after configuration;
//   - option_targets.go owns public target-list options;
//   - target_list.go and target_error.go own target-list validation and copying;
//   - runner.go owns the Runner type;
//   - runner_constructor.go owns Runner construction;
//   - runner_run.go owns schedule-loop lifecycle;
//   - runner_cycle.go owns one probe cycle and evaluator-error normalization;
//   - runner_snapshot.go owns public read methods;
//   - runner_error.go owns runner error sentinels;
//   - snapshot.go owns Snapshot invariants and predicates;
//   - snapshot_clone.go owns defensive copy helpers;
//   - observation.go owns the private stored observation payload;
//   - stale.go owns stale calculation;
//   - store.go owns the private latest-snapshot cache.
//
// # Dependency policy
//
// Production code depends only on the Go standard library, arcoris.dev/health,
// arcoris.dev/chrono/clock, arcoris.dev/chrono/delay, and
// arcoris.dev/snapshot.
package probe
