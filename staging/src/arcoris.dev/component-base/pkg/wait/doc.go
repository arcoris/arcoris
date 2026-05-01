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

// Package wait provides low-level runtime wait mechanics for ARCORIS component
// internals.
//
// The package belongs to component-base because it defines shared mechanical
// primitives used by component loops, shutdown paths, timers, and polling code.
// It owns waiting mechanics, not runtime policy. Higher-level packages decide
// what should be retried, scheduled, queued, logged, traced, exported, or treated
// as domain state.
//
// # Scope
//
// This package owns:
//
//   - condition contracts and deterministic condition composition;
//   - cancellable one-shot runtime delays;
//   - owner-controlled real-runtime timers;
//   - wait-owned interruption and timeout classification;
//   - bounded positive jitter for real-runtime desynchronization;
//   - functional options for mechanical wait behavior;
//   - fixed-interval condition loops.
//
// It intentionally stays small. It does not provide retry policy, backoff
// policy, scheduler policy, queue policy, worker lifecycle policy, admission
// policy, logging, tracing, metrics, distributed ordering, or domain-specific
// behavior.
//
// # Condition model
//
// A ConditionFunc reports exactly one of three outcomes:
//
//   - done=true, err=nil means the wait completed successfully;
//   - done=false, err=nil means the wait should continue;
//   - err!=nil means the wait must stop and return that condition-owned error.
//
// Wait loops must return condition-owned errors unchanged. They must not
// reinterpret, wrap, aggregate, retry, suppress, or convert errors returned by a
// condition. In particular, raw context.Canceled and context.DeadlineExceeded
// values returned by a condition remain condition-owned errors and are not
// classified as wait-owned interruption or timeout.
//
// Condition helpers evaluate composed conditions sequentially and
// deterministically. They do not start goroutines, evaluate conditions
// concurrently, recover panics, or add synchronization around caller-owned code.
//
// # Context ownership
//
// Public wait primitives own the context they receive for their mechanical
// operation. They validate that the context is non-nil at the API boundary.
//
// Conditions receive that same context. Conditions should observe the context
// when evaluation can block, perform I/O, acquire resources, or depend on
// cancellation or deadlines. Pure in-memory predicates may ignore the context.
//
// # Wait-owned errors
//
// Context stops observed by wait primitives are wait-owned errors. Cancellation
// is classified as ErrInterrupted. Deadline expiration is classified as
// ErrTimeout, which is also a specialized ErrInterrupted.
//
// Wait-owned context errors preserve context causes. Ordinary cancellation and
// deadline expiration match context.Canceled or context.DeadlineExceeded through
// errors.Is. Cancellation or deadline with a custom cause preserves that cause
// and also preserves matching against the standard context sentinel when the
// context exposes one.
//
// Wait-owned classification is added only by wait primitives that own the
// context stop. Raw context errors returned by conditions remain condition-owned
// and must not be reclassified by the wait loop.
//
// # Timer model
//
// Timer is an owner-controlled one-shot real-runtime timer backed by time.Timer.
// It is useful when code needs explicit stop, drain, reset, channel access, or
// context-aware waiting with wait-owned error classification.
//
// Timer values must be created with NewTimer. The zero value is invalid. Timer
// follows a single-owner coordination model and must not be copied after
// construction. Callers that share Stop, StopAndDrain, Reset, Wait, or channel
// receives across goroutines must provide their own synchronization.
//
// Timer is not a replacement for clock.Timer. Components that need deterministic
// fake-time tests or a runtime clock abstraction should use package
// arcoris.dev/component-base/pkg/clock. Components that only need simple real
// runtime waiting may use Delay, Timer, and Until from this package.
//
// # Relationship to package clock
//
// Package clock abstracts time sources and provides real and fake clocks,
// timers, and tickers. Package wait uses real runtime time directly for its
// mechanical delays and timers. It does not accept or store a clock.Clock
// because deterministic fake time is outside this package's responsibility.
//
// Code that must be tested without real time should build on package clock. Code
// that needs only a small real-runtime wait layer can use package wait.
//
// # Jitter model
//
// Jitter is positive and one-sided. It may extend a base duration by a bounded
// random delta, but it never shortens the base duration. This keeps jitter safe
// for conservative intervals where waking earlier than the base cadence would
// violate caller policy.
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
// # File ownership
//
// The package keeps files split by responsibility:
//
//   - condition.go owns ConditionFunc;
//   - condition_helpers.go owns condition helper composition;
//   - error_kind.go, error_interrupted.go, error_timeout.go, and
//     error_context_stop.go own wait error classification;
//   - delay.go owns cancellable one-shot delays;
//   - timer.go owns the owner-controlled runtime timer;
//   - until.go owns fixed-interval condition loops;
//   - jitter.go owns positive one-sided jitter;
//   - options.go owns mechanical functional options;
//   - validate.go owns package-local validation helpers;
//   - nocopy.go owns the static-analysis copy marker.
//
// # Nil input and panic policy
//
// Nil contexts, nil predicates, nil conditions, nil timers, nil options, and
// invalid jitter factors are programming errors. Public wait primitives and
// helper constructors panic immediately when they receive invalid inputs instead
// of returning values that fail later inside a runtime loop.
//
// Panics raised by condition implementations are not recovered by this package.
// Panic recovery, if required, belongs to the loop owner or to an explicit
// wrapper in a higher-level package.
//
// # Non-goals
//
// This package does not provide:
//
//   - retry policy or retry budgets;
//   - backoff algorithms;
//   - scheduler policy;
//   - queue policy;
//   - admission policy;
//   - worker lifecycle policy;
//   - lease ownership or fencing;
//   - distributed coordination or ordering;
//   - logging, tracing, metrics, or exporters;
//   - deterministic fake time.
//
// Higher-level packages may build those behaviors using ConditionFunc, Delay,
// Timer, Jitter, Until, and wait-owned errors, but this package must remain a
// small mechanical runtime wait layer.
//
// # Dependency policy
//
// Production code in this package depends only on the Go standard library.
package wait
