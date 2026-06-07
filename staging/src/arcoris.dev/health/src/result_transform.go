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

package health

import "time"

// WithCause returns a copy of r with cause attached as the internal lower-level
// cause.
//
// Cause is preserved for logs, tests, diagnostics, and error classification. It
// MUST NOT be exposed by public transport adapters by default.
func (r Result) WithCause(cause error) Result {
	r.Cause = cause
	return r
}

// WithObserved returns a copy of r with observed set as the observation time.
//
// Observed should represent when the health state was produced or normalized,
// not necessarily when it is later rendered by an adapter.
func (r Result) WithObserved(observed time.Time) Result {
	r.Observed = observed
	return r
}

// WithDuration returns a copy of r with duration set as the observation
// duration.
//
// The method does not reject negative values because Result is a plain value
// type. Call Normalize to defensively repair invalid durations at ownership
// boundaries.
func (r Result) WithDuration(d time.Duration) Result {
	r.Duration = d
	return r
}

// WithMessage returns a copy of r with message set as the safe human-readable
// message.
//
// The caller owns message safety. Message MUST remain suitable for the adapters
// and diagnostics that will render it.
func (r Result) WithMessage(message string) Result {
	r.Message = message
	return r
}

// WithReason returns a copy of r with reason set as the machine-readable reason.
//
// Reason should be stable enough for policy, diagnostics, and tests. It should
// not contain caller-specific details, timestamps, resource identifiers, or raw
// low-level error strings.
func (r Result) WithReason(reason Reason) Result {
	r.Reason = reason
	return r
}

// Normalize returns a defensively normalized copy of r.
//
// Normalize is intended for checker, evaluator, resolver, and report boundaries
// where ownership of the result is known. It fills an empty Name with
// defaultName, replaces invalid statuses with StatusUnknown, fills a zero
// Observed value with observed, and clamps negative Duration to zero.
//
// Normalize does not rewrite Reason, Message, or Cause. Those fields preserve
// checker-owned semantics and must be interpreted by the caller.
func (r Result) Normalize(defaultName string, observed time.Time) Result {
	if r.Name == "" {
		r.Name = defaultName
	}
	if !r.Status.IsValid() {
		r.Status = StatusUnknown
	}
	if r.Observed.IsZero() {
		r.Observed = observed
	}
	if r.Duration < 0 {
		r.Duration = 0
	}

	return r
}
