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

package retry

import "time"

// Outcome describes how one retry execution completed.
//
// Outcome is immutable completion metadata. It records how many operation calls
// were made, when the retry execution started and finished, which operation-owned
// error was observed last, and why the retry execution stopped.
//
// Outcome does not execute operations, classify errors, compute retry delays,
// wait on timers, observe contexts, notify observers, or create retry-owned error
// wrappers. Those responsibilities belong to the retry loop, classifier,
// delay, clock, observer, and error layers.
//
// Outcome values are copyable. They must not contain locks, channels, mutable
// delay sequences, observers, operations, timers, or references back to the
// retry execution that produced them.
//
// The zero Outcome is invalid. This prevents empty metadata from accidentally
// looking like a successful execution.
type Outcome struct {
	// Attempts is the number of operation calls performed by the retry execution.
	//
	// Attempts is zero only when retry execution stopped before the first
	// operation call, such as when the retry-owned context was already stopped.
	// Successful, non-retryable, exhausted, and post-attempt interrupted outcomes
	// must have a non-zero attempt count.
	Attempts uint

	// StartedAt is the retry clock time at which retry execution began.
	//
	// StartedAt is assigned once per Do or DoValue execution before the first
	// attempt is considered. It is retry execution metadata, not attempt metadata.
	// Individual operation calls are described by Attempt.StartedAt.
	StartedAt time.Time

	// FinishedAt is the retry clock time at which retry execution reached its
	// terminal decision.
	//
	// FinishedAt must not be earlier than StartedAt for a valid Outcome. It is
	// intended for diagnostics, observer events, elapsed-time accounting, and
	// tests. It is not a distributed ordering source.
	FinishedAt time.Time

	// LastErr is the last operation-owned error observed by the retry execution.
	//
	// LastErr is nil for successful outcomes. It is also nil when retry execution
	// is interrupted before the first operation call. For failed outcomes that
	// happen after an operation attempt, LastErr should preserve the operation
	// error that led retry to stop, exhaust limits, or wait for another attempt.
	LastErr error

	// Reason explains why retry execution stopped.
	//
	// Reason classifies the terminal retry decision. It does not carry the
	// operation error, context cause, delay sequence state, or observer state.
	Reason StopReason
}

// IsZero reports whether o is the zero Outcome value.
//
// The zero Outcome is not valid completion metadata. It is useful as an omitted
// value in events that are not terminal retry-stop events.
func (o Outcome) IsZero() bool {
	return o.Attempts == 0 &&
		o.StartedAt.IsZero() &&
		o.FinishedAt.IsZero() &&
		o.LastErr == nil &&
		o.Reason == 0
}

// IsValid reports whether o is internally consistent retry completion metadata.
//
// IsValid checks structural invariants only. It does not check whether timestamps
// came from a particular clock, whether LastErr is retryable, whether Attempts
// matches a specific delay sequence, or whether Reason was produced by a real
// retry loop.
func (o Outcome) IsValid() bool {
	if !o.Reason.IsValid() {
		return false
	}
	if o.StartedAt.IsZero() || o.FinishedAt.IsZero() {
		return false
	}
	if o.FinishedAt.Before(o.StartedAt) {
		return false
	}

	switch o.Reason {
	case StopReasonSucceeded:
		return o.Attempts != 0 && o.LastErr == nil

	case StopReasonNonRetryable,
		StopReasonMaxAttempts,
		StopReasonMaxElapsed,
		StopReasonDeadline,
		StopReasonDelayExhausted:
		return o.Attempts != 0 && o.LastErr != nil

	case StopReasonInterrupted:
		if o.Attempts == 0 {
			return o.LastErr == nil
		}
		return o.LastErr != nil

	default:
		return false
	}
}

// Succeeded reports whether o records a successful retry execution.
//
// Succeeded returns true only for valid successful outcomes. Malformed metadata
// with Reason set to StopReasonSucceeded still returns false when the rest of the
// outcome is not internally consistent.
func (o Outcome) Succeeded() bool {
	return o.IsValid() && o.Reason.Succeeded()
}

// Failed reports whether o records an unsuccessful retry execution.
//
// Failed returns true only for valid failed outcomes. Invalid metadata is not
// treated as an ordinary retry failure.
func (o Outcome) Failed() bool {
	return o.IsValid() && o.Reason.Failed()
}

// Exhausted reports whether o records retry-owned exhaustion.
//
// Exhaustion includes max-attempt exhaustion, max-elapsed exhaustion, context
// deadline budget exhaustion, and finite delay sequence exhaustion.
// Non-retryable operation errors and retry-owned interruptions are failures, but
// they are not exhaustion.
func (o Outcome) Exhausted() bool {
	return o.IsValid() && o.Reason.Exhausted()
}

// Interrupted reports whether o records retry-owned interruption.
//
// Interruption means retry execution stopped because retry observed its owning
// context stop at a retry boundary, such as before an attempt or while waiting
// between attempts.
func (o Outcome) Interrupted() bool {
	return o.IsValid() && o.Reason.Interrupted()
}

// Duration returns the elapsed retry execution duration.
//
// Duration returns zero when StartedAt or FinishedAt is zero, or when FinishedAt
// is earlier than StartedAt. Valid outcomes therefore always return a
// non-negative duration.
func (o Outcome) Duration() time.Duration {
	if o.StartedAt.IsZero() || o.FinishedAt.IsZero() {
		return 0
	}
	if o.FinishedAt.Before(o.StartedAt) {
		return 0
	}

	return o.FinishedAt.Sub(o.StartedAt)
}
