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

// Package run provides context-first task orchestration for ARCORIS component
// internals.
//
// The package belongs to arcoris.dev/runtime because component runtimes often
// need to start a small set of related goroutines under one context, cancel
// siblings after a task failure, wait for deterministic completion, and return
// useful named task errors. It owns those task-group mechanics only. It does not own
// process signals, lifecycle states, supervision, restart policy, retry policy,
// logging, metrics, tracing, process exit, worker pools, or job queues.
//
// # Scope
//
// This package owns:
//
//   - the Task function contract;
//   - owner-created Group values;
//   - shared group context cancellation;
//   - named goroutine startup through Group.Go;
//   - deterministic waiting through Group.Wait;
//   - named TaskError values for task failures;
//   - context-first helper functions for common task returns.
//
// It intentionally stays mechanical. Higher-level packages decide what a task
// means, how many tasks to start, whether a context comes from a process signal
// or another owner, how failures should affect lifecycle state, and how errors
// are logged, exported, retried, or reported.
//
// # Task model
//
// Task is a context-aware unit of work:
//
//	group := run.NewGroup(ctx)
//	group.Go("dispatcher", dispatcher.Run)
//	group.Go("controller", controller.Run)
//	err := group.Wait()
//
// A Task receives the Group context. It should observe that context whenever it
// blocks, loops, waits on I/O, sends, receives, sleeps, or owns resources that
// must be released during shutdown. Returning nil means the task completed
// successfully or stopped cleanly. Returning a non-nil error means the task
// failed and should be reported through Wait.
//
// Package run does not recover panics from Task implementations. Panic recovery,
// if required, is owner policy and should be implemented by explicitly wrapping
// tasks before passing them to Group.Go.
//
// # Group model
//
// Group is created with NewGroup. The zero Group is invalid and a Group must not
// be copied after construction. NewGroup derives a child context from the parent
// using context.WithCancelCause.
//
// Go starts a named task in a goroutine. Task names must be non-empty, trimmed,
// and unique within the Group. Go may be called concurrently while the Group is
// open. Wait, Cancel, and fail-fast task errors close the Group for new
// submissions.
//
// Wait is the join point. It closes the Group for new submissions, waits for all
// started tasks, cancels the Group context after completion, builds the
// configured task error result, and caches that result. Later Wait calls return
// the same cached error.
//
// # Error model
//
// A non-nil task return is recorded as TaskError{Name, Err}. TaskError matches
// ErrTaskFailed with errors.Is and unwraps the original task error, so
// errors.Is(waitErr, originalErr) and errors.As(waitErr, *TaskError) work on
// Wait results.
//
// ErrorModeJoin returns all recorded task errors through errors.Join, ordered by
// task submission sequence rather than goroutine completion order. ErrorModeFirst
// returns the first observed task error. TaskErrors walks both ordinary unwrap
// chains and joined-error trees to recover TaskError values from arbitrary error
// results.
//
// # Cancellation model
//
// Groups cancel on task errors by default. WithCancelOnError(false) records task
// errors without cancelling the shared context. Owner cancellation through Cancel
// closes the Group for new submissions and cancels the context, but it does not
// create a task error.
//
// Cancel(nil) uses context.Canceled. Cancel with a custom cause preserves that
// cause according to context.WithCancelCause semantics. If a task error wins
// fail-fast cancellation, the Group context cause is the TaskError for that
// first recorded task failure.
//
// # Helper functions
//
// Wait(ctx) blocks until ctx is done and returns context.Cause(ctx), falling
// back to ctx.Err when no explicit cause is available. IsContextStop,
// IgnoreContextCanceled, and IgnoreContextStop help tasks distinguish expected
// context shutdown from unrelated work errors without importing another
// adjacent ARCORIS package.
//
// # Relationship to other adjacent ARCORIS packages
//
// Package signals owns process signal registration and signal-derived contexts.
// Package run does not import signals and does not decide which process signals
// should cancel a Group.
//
// Package lifecycle owns component lifecycle states and transitions. Package run
// does not import lifecycle and does not mark components starting, running,
// stopping, stopped, or failed.
//
// Package wait owns low-level wait mechanics. Package run provides a few
// context-first helpers because tasks commonly need them, but it does not become
// a timer, retry, backoff, polling, or scheduler package.
//
// # Non-goals
//
// This package does not provide:
//
//   - actors or actor groups;
//   - worker pools, job queues, Submit, or TrySubmit APIs;
//   - concurrency limits;
//   - supervisors, restarts, retries, or backoff;
//   - process signal handling;
//   - lifecycle transitions;
//   - logging, metrics, tracing, or hooks;
//   - panic recovery by default;
//   - process exit behavior.
//
// # File ownership
//
// The package keeps files split by responsibility:
//
//   - task.go owns the Task contract;
//   - group.go owns Group construction and public methods;
//   - group_option.go owns GroupOption, groupConfig, ErrorMode, and option
//     normalization;
//   - group_state.go owns internal Group state transitions and submission
//     guards;
//   - group_error.go owns internal task error recording, ordering, and Wait
//     error construction;
//   - error.go owns ErrTaskFailed, TaskError, and TaskErrors;
//   - wait.go owns the context-blocking sentinel task helper;
//   - context_stop.go owns context-stop classification and ignore helpers;
//   - context_cause.go owns package-local context cause fallback behavior;
//   - validate.go owns package-local validation helpers and stable panic
//     messages;
//   - nocopy.go owns the static-analysis copy marker.
//
// # Dependency policy
//
// Production code in this package depends only on the Go standard library.
package run
