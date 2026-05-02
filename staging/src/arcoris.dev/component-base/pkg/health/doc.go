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

// Package health provides local health check contracts and synchronous
// evaluation mechanics for ARCORIS component internals.
//
// # Package scope
//
// The package owns local, in-process health primitives:
//
//   - Checker and CheckFunc contracts;
//   - Status, Reason, Target, Result, and Report values;
//   - Registry storage for target-scoped checks;
//   - Gate for owner-published mutable health state;
//   - Evaluator for synchronous target evaluation;
//   - evaluator-owned timeout, cancellation, and panic normalization.
//
// It does not expose HTTP handlers, implement Kubernetes probes, publish
// metrics, perform service discovery, run distributed health consensus, drive
// lifecycle transitions, or make restart, admission, routing, scheduling,
// logging, tracing, or alerting decisions.
//
// # Health check contract
//
// A Checker owns one stable check name and returns one Result for each Check
// call. Checkers should observe the supplied context when work can block, wait
// on I/O, acquire resources, or depend on cancellation. Pure in-memory checks
// may ignore context.
//
// Checkers should return Results instead of panicking. Evaluator recovers checker
// panics so aggregation remains robust, but a panic still represents a checker
// programming fault and is converted into an unhealthy Result with ReasonPanic.
//
// # Result model
//
// Status is the primary health state. Reason is a stable machine-readable cause.
// Message is the safe human-readable text that adapters may expose. Cause is
// internal diagnostic detail and must not be exposed by public adapters by
// default.
//
// Target identifies the scope being evaluated: startup, live, or ready. Target
// policies decide which statuses pass for a target. The core package keeps these
// policies transport-neutral and does not map them to HTTP status codes or gRPC
// serving states.
//
// # Reason taxonomy
//
// Built-in Reasons cover common operational categories that are broadly useful
// across distributed runtime components:
//
//   - observation and evaluator boundaries, such as missing observations,
//     timeouts, cancellations, and recovered panics;
//   - lifecycle and operation phases, such as startup, draining, and shutdown;
//   - dependency availability and dependency degradation;
//   - load, backpressure, rate limiting, admission closure, capacity exhaustion,
//     and resource exhaustion;
//   - freshness and synchronization state, such as stale, not-synced, failed
//     sync, and lagging views;
//   - distributed connectivity partitioning;
//   - configuration errors and fatal owner-defined failures.
//
// Reason is deliberately extensible. Domain packages may define their own valid
// lower_snake_case reasons when the core taxonomy is too coarse. Core code should
// not grow storage-specific, scheduler-specific, security-specific,
// lease-specific, or transport-specific reasons unless those causes become
// common component-base concepts. Dynamic details belong in Message when safe,
// or in Cause when internal.
//
// # Evaluator ownership
//
// Evaluator reads checks from Registry at evaluation time, executes them
// sequentially in registration order, normalizes individual Results, aggregates
// the most severe Status, and returns a Report. It does not retain Reports, run
// periodic probes, start background goroutines outside a single timeout boundary,
// retry checks, or cache results.
//
// # Timeout and context model
//
// Evaluator gives each check a cooperative context. With a positive timeout, it
// uses a per-check context timeout and returns an unknown timeout Result if that
// context finishes first. A checker that ignores context may outlive the
// evaluator result; health only owns the caller-visible evaluation boundary, not
// forced goroutine termination.
//
// Parent cancellation is reported as an unknown canceled Result. Custom context
// causes remain attached to Result.Cause through the context package.
//
// # Clock model
//
// Evaluator depends on clock.PassiveClock for observation timestamps and elapsed
// duration measurement. It does not create timers, tickers, sleeps, or scheduler
// loops. Timeout enforcement uses context deadlines because it is an
// evaluator-owned execution boundary, not a reusable clock scheduler.
//
// # Registry, gate, and policy
//
// Registry stores target-scoped checks in deterministic registration order. Gate
// is a concurrency-safe Checker for owner-published health state. TargetPolicy
// interprets a Status for one target without prescribing transport or runtime
// action.
//
// # File ownership
//
//   - check.go owns Checker and check name validation;
//   - check_func.go owns function-backed checker adapters;
//   - identifier.go owns shared lower_snake_case identifier syntax;
//   - status.go owns Status values and status ordering;
//   - reason.go owns Reason values and reason classification;
//   - target.go owns Target values and target enumeration;
//   - result.go owns single-check Result values;
//   - report.go owns target-level Report values;
//   - policy.go owns target status policy;
//   - registry.go owns target-scoped check registration;
//   - registry_error.go owns registry error sentinels and typed errors;
//   - registry_validate.go owns registration batch validation;
//   - gate.go owns mutable owner-published checker state;
//   - shutdown.go owns shutdown and drain check adapters;
//   - evaluator.go owns Evaluator construction and public evaluation;
//   - evaluator_run.go owns single-check execution and normalization;
//   - evaluator_error.go owns evaluator error sentinels;
//   - evaluator_panic.go owns panic cause preservation;
//   - evaluator_option*.go files own evaluator option domains.
//
// # Dependency policy
//
// Production code depends only on the Go standard library and package clock.
package health
