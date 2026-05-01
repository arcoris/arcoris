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

// Package wait provides low-level wait condition contracts, owner-controlled
// runtime timers, cancellable delay primitives, positive-jitter helpers, and
// small runtime wait loops for ARCORIS component code.
//
// The package is part of component-base. It defines the shared wait vocabulary
// used by polling, retry, backoff, controller, queue, worker, shutdown, and
// reconciliation code without owning those higher-level policies itself.
//
// # Scope
//
// This package owns the low-level wait mechanics:
//
//   - the ConditionFunc contract;
//   - deterministic helper conditions;
//   - adaptation from context-aware boolean predicates;
//   - sequential condition composition;
//   - wait-owned interruption and timeout classification;
//   - owner-controlled one-shot runtime timers;
//   - cancellable one-shot delays;
//   - bounded positive jitter for runtime desynchronization;
//   - simple fixed-interval condition loops.
//
// The package intentionally does not own higher-level runtime policy. Retry
// limits, backoff curves, scheduler policy, queue policy, admission policy,
// worker lifecycle policy, metrics, tracing, logging, and observability belong
// to packages that build on top of these primitives.
//
// # Condition outcomes
//
// A ConditionFunc reports one of three outcomes:
//
//   - done=true, err=nil means the wait completed successfully;
//   - done=false, err=nil means the wait should continue;
//   - err!=nil means the wait must stop and return that error.
//
// Condition errors are terminal condition-owned results. Wait loops return them
// unchanged and do not reinterpret, wrap, aggregate, retry, suppress, or convert
// them.
//
// # Context ownership
//
// Conditions are context-aware by default. The context is owned by the wait
// operation that evaluates the condition or performs the delay.
//
// Condition implementations SHOULD observe the context when evaluation can
// block, perform I/O, acquire resources, wait for another component, or depend
// on cancellation or deadlines. Pure in-memory predicates MAY ignore the
// context.
//
// Public wait primitives reject nil contexts at their API boundary. Use
// context.Background when no narrower cancellation scope is available. Helper
// functions pass the received context through unchanged; nil-context validation
// does not belong inside every condition helper layer.
//
// # Wait-owned errors
//
// Context stops observed by wait primitives are wait-owned errors. Cancellation
// is classified as ErrInterrupted. Deadline expiration is classified as
// ErrTimeout, which is also a specialized ErrInterrupted.
//
// Raw context.Canceled and context.DeadlineExceeded values returned by a
// condition remain condition-owned errors and are not automatically classified
// as wait-owned interruption or timeout. Only the wait primitive that owns the
// context may add wait-owned classification.
//
// # Evaluation model
//
// Wait loops may evaluate a condition repeatedly. Condition implementations MUST
// be safe to call repeatedly by the wait operation that owns them.
//
// Conditions MUST NOT rely on an exact number of evaluations. A wait may stop
// early due to success, error, context cancellation, timeout, shutdown, or a
// runtime loop policy. A wait implementation may also perform an immediate
// first evaluation before sleeping, depending on its public contract.
//
// Condition helpers in this package evaluate composed conditions
// deterministically and sequentially. They do not start goroutines, evaluate
// conditions concurrently, recover panics, or add synchronization around caller
// code.
//
// # Timer model
//
// Timer is an owner-controlled one-shot runtime timer. It is useful when a
// caller needs explicit lifecycle control over a real-time wait: stop, drain,
// reset, channel access, or context-aware waiting with wait-owned error
// classification.
//
// Timer values must be created with NewTimer and must not be copied after
// construction. The zero value is not usable because it has no underlying
// runtime timer. Timer follows a single-owner coordination model; components
// that share Stop, Reset, or Wait across goroutines must provide their own
// synchronization.
//
// # Jitter model
//
// Jitter is positive and one-sided: it may extend a base duration by a bounded
// random delta, but it never shortens the base duration. This makes jitter safe
// for conservative wait intervals where waking up earlier than the base cadence
// would violate a caller's policy.
//
// Jitter is only a duration-spreading primitive. Backoff growth, retry budgets,
// retryability rules, fairness, and SLO policy belong to higher-level packages.
//
// # Options model
//
// Option values configure narrow mechanical behavior for wait primitives without
// exposing mutable configuration structs. Options are applied in the order they
// are supplied. When more than one option configures the same domain, the later
// option wins.
//
// Options must remain low-level. They may configure mechanics such as interval
// jitter, but they must not encode retryability rules, retry budgets, scheduler
// policy, queue policy, observability policy, or domain-specific decisions.
//
// # Nil input and panic policy
//
// Nil contexts, nil predicates, nil conditions, nil timers, nil options, and
// invalid jitter factors are programming errors. Public wait primitives and helper
// constructors panic immediately when they receive invalid inputs instead of
// returning values that fail later inside a runtime loop.
//
// Panics raised by condition implementations are not recovered by this package.
// Panic recovery, if required by a runtime owner, should be implemented by the
// loop owner or by an explicit wrapper in a higher-level package.
//
// # Non-goals
//
// This package does not define:
//
//   - retry policy;
//   - backoff algorithms;
//   - scheduler policy;
//   - queue policy;
//   - admission policy;
//   - worker lifecycle policy;
//   - distributed coordination;
//   - metrics or observability instruments.
//
// Higher-level packages may build those behaviors using ConditionFunc, Timer,
// Delay, Jitter, Until, and wait-owned errors, but this package must remain a
// small runtime wait layer.
//
// # Dependency policy
//
// Production code in this package should depend only on the Go standard library.
package wait

import "context"

// ConditionFunc evaluates whether a wait condition has been satisfied.
//
// A condition returns done=true when the caller should stop waiting
// successfully. It returns done=false with a nil error when the caller should
// keep waiting. It returns a non-nil error when waiting must stop and the error
// must be returned to the wait owner.
//
// The context is owned by the wait operation. Condition implementations SHOULD
// observe ctx when evaluation can block, perform I/O, acquire resources, or
// depend on cancellation or deadlines. Pure in-memory conditions MAY ignore ctx.
//
// ConditionFunc implementations MUST be safe to call repeatedly by the wait loop
// that owns them. They MUST NOT assume an exact number of evaluations because a
// wait may stop early due to success, error, context cancellation, timeout, or
// shutdown.
//
// Callers MUST pass a non-nil context. Use context.Background when no specific
// cancellation scope is available.
type ConditionFunc func(ctx context.Context) (done bool, err error)
