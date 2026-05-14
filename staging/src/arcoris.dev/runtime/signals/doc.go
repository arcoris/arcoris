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

// Package signals provides process signal integration primitives for ARCORIS
// component runtimes.
//
// The package belongs to arcoris.dev/runtime because long-running ARCORIS components,
// command entrypoints, agents, workers, controllers, examples, and test harnesses
// need one stable way to bridge process-level OS signals into ordinary runtime
// coordination primitives. It owns process signal registration, signal set
// normalization, signal-to-context cancellation, and graceful-shutdown triggers.
// It does not own component lifecycle, supervision, restart policy, logging,
// metrics, tracing, process exit, or application shutdown deadlines.
//
// # Scope
//
// This package owns:
//
//   - platform-aware signal sets for shutdown, reload, and diagnostics;
//   - copyable signal Event values;
//   - typed SignalError cancellation causes;
//   - owner-controlled signal Subscription values;
//   - context cancellation when a configured signal is received;
//   - graceful shutdown coordination with optional repeated-signal escalation;
//   - test seams around os/signal registration.
//
// It intentionally stays mechanical. Higher-level packages decide what a signal
// means for lifecycle transitions, worker shutdown, retry cancellation,
// configuration reload, diagnostics, logging, metrics, tracing, health, or
// process exit.
//
// # Signal sets
//
// ShutdownSignals returns the project-wide default signal set that should
// initiate graceful shutdown for ordinary ARCORIS processes. ReloadSignals and
// DiagnosticSignals return platform-specific sets for optional owners that
// deliberately support those behaviors.
//
// Signal set helpers such as Clone, Unique, Merge, and Contains are pure value
// helpers. They preserve caller order where order matters and reject nil signals
// as programming errors. Default signal set functions return fresh slices so
// callers cannot mutate package-owned state.
//
// # Subscription model
//
// Subscription is the low-level owner-controlled registration primitive. It
// registers one channel with package os/signal and exposes C, Wait, Stop, and
// Done operations. Callers that create a Subscription must call Stop when the
// subscription is no longer needed.
//
// C and Wait receive from the same channel. A direct receive from C consumes the
// same signal value that Wait would otherwise observe, so owners should choose a
// single receive path or deliberately coordinate any competition between those
// APIs.
//
// Subscription does not close the signal delivery channel. The channel is owned
// by the subscription, but os/signal delivery and user receives may race with
// shutdown. Stop unregisters the channel from os/signal and closes Done so
// waiting goroutines can leave deterministically.
//
// Subscription values must be created with Subscribe or SubscribeWithOptions and
// must not be copied after first use.
//
// # Context model
//
// NotifyContext derives a child context that is cancelled when the parent stops,
// when StopFunc is called, or when one configured signal is received. Signal
// cancellation uses SignalError as the context cancellation cause. Parent
// cancellation preserves the parent cause when available.
//
// StopFunc must be called by the owner. It unregisters signal delivery, releases
// resources, and cancels the returned context if it has not already been
// cancelled.
//
// # Shutdown controller model
//
// ShutdownController is the higher-level graceful shutdown primitive. The first
// configured shutdown signal records the first Event and cancels the controller
// context with SignalError. Shutdown signals are registered when the controller
// is constructed. Escalation signals are registered only after shutdown starts,
// which avoids intercepting escalation-only process signals before the
// controller owns shutdown. Repeated escalation signals are delivered through a
// best-effort escalation channel when escalation is enabled.
//
// ShutdownController never calls os.Exit, never panics on repeated signals, and
// never drives lifecycle transitions directly. Applications may decide that a
// repeated signal should shorten a grace period, dump diagnostics, return from
// main, or exit the process. Escalation delivery is best-effort and does not
// exit the process. That policy belongs to the application owner, not to
// arcoris.dev/runtime.
//
// # Error and cause model
//
// SignalError is the typed error used when a signal cancels a context. It matches
// ErrSignal with errors.Is and exposes the received Event. Cause extracts a
// signal Event from context.Cause when a context was cancelled by this package.
//
// ErrStopped reports owner-initiated subscription shutdown to Wait callers. It
// does not mean an OS signal was received.
//
// # Concurrency and cleanup
//
// Subscription and ShutdownController methods are safe for concurrent use unless
// a method comment says otherwise. Stop operations are idempotent. Owners should
// still treat Subscription and ShutdownController values as single runtime
// owners and should not copy them after construction.
//
// # Relationship to lifecycle
//
// Package lifecycle owns component lifecycle states, transition validation,
// snapshots, guards, observers, and lifecycle waiters. Package signals only
// provides process signal triggers. A component may choose to call lifecycle
// transition methods when a signal-derived context is cancelled, but this
// package does not import or drive lifecycle directly.
//
// # Relationship to wait
//
// Package wait owns low-level waiting mechanics and wait-owned context error
// classification. Signal-derived contexts can be passed to wait primitives. When
// a signal cancels such a context, wait primitives may wrap or classify the
// context stop while preserving the signal cause according to their own error
// model.
//
// # Non-goals
//
// This package does not provide:
//
//   - process supervision;
//   - restart policy;
//   - lifecycle state transitions;
//   - graceful shutdown deadlines or budgets;
//   - forced process exit;
//   - panic-on-second-signal behavior;
//   - retry, backoff, queue, scheduler, or admission policy;
//   - logging, metrics, tracing, or exporters;
//   - CLI command routing;
//   - daemonization or systemd integration;
//   - distributed coordination or leader election.
//
// # File ownership
//
// The package keeps files split by responsibility:
//
//   - event.go owns Event;
//   - error.go owns ErrSignal, ErrStopped, SignalError, and Cause;
//   - set.go owns signal set value helpers;
//   - set_unix.go, set_windows.go, and set_other.go own platform sets;
//   - validate.go owns package-local validation helpers;
//   - notifier.go and os_notifier.go own the package-local os/signal seam;
//   - subscription_config.go owns Subscription construction config;
//   - option_subscription.go owns Subscription construction options;
//   - subscription.go owns signal registration lifecycle;
//   - context.go owns NotifyContext;
//   - shutdown_config.go owns ShutdownController construction config;
//   - option_shutdown.go owns ShutdownController construction options;
//   - shutdown.go owns graceful shutdown coordination;
//   - nocopy.go owns the static-analysis copy marker.
//
// # Dependency policy
//
// Production code in this package depends only on the Go standard library.
package signals
