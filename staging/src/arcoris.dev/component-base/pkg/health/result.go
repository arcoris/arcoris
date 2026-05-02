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

package health

import "time"

// Result describes one health observation produced by a checker, gate, cached
// probe, or another health source.
//
// Result is intentionally transport-neutral. It does not define HTTP response
// codes, gRPC serving states, JSON rendering, metric labels, restart decisions,
// readiness outcomes, admission behavior, or scheduler policy. Those concerns
// belong to adapter packages and higher-level runtime owners.
//
// The zero value is an unnamed StatusUnknown result. It is safe to allocate in
// larger structs before the first health observation is available. Callers that
// publish, aggregate, or expose results SHOULD normalize them at the boundary
// where checker ownership and observation time are known.
//
// Name identifies the logical check that produced the result. A checker-owned
// result SHOULD set Name to the checker name. Aggregators MAY fill an empty Name
// from the checker that returned the result.
//
// Status is the primary operational health state.
//
// Reason is a stable, machine-readable explanation for Status when a reason is
// useful. A healthy result MAY leave Reason empty. Non-healthy results SHOULD
// provide a reason when the owner can classify the condition without leaking
// private implementation detail.
//
// Message is a safe, short, human-readable explanation. Message MUST NOT contain
// private operational data. Transport adapters may expose Message directly.
//
// Observed is the time at which the result was observed or normalized. A zero
// value means the result has not been timestamped yet.
//
// Duration is the amount of time spent producing the observation. A zero value is
// valid. Negative durations are invalid and SHOULD be normalized before
// aggregation or exposure.
//
// Cause preserves the internal lower-level failure cause. Cause is intentionally
// not a public diagnostic field. Transport adapters MUST NOT expose Cause by
// default. Logs, tests, and owner-controlled diagnostics may inspect Cause when
// they have permission to handle internal details.
type Result struct {
	Name     string
	Status   Status
	Reason   Reason
	Message  string
	Observed time.Time
	Duration time.Duration
	Cause    error
}

// Healthy returns a healthy result for name.
//
// A healthy result intentionally has no reason or message by default. Callers
// that need a human-readable explanation for diagnostics may set Message with a
// follow-up value transformation, but most healthy results should remain compact.
func Healthy(name string) Result {
	return Result{
		Name:   name,
		Status: StatusHealthy,
	}
}

// Starting returns a starting result for name.
//
// Starting means the checked component or subsystem is still bootstrapping. It
// is not a terminal failure. Target policies decide whether starting is accepted
// for a specific health target.
func Starting(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusStarting,
		Reason:  reason,
		Message: message,
	}
}

// Degraded returns a degraded result for name.
//
// Degraded means the checked component or subsystem still has usable capability,
// but is operating with reduced capability, reduced confidence, or active
// protective behavior. Degraded MUST remain distinct from unhealthy so runtime
// owners can make target-specific decisions.
func Degraded(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusDegraded,
		Reason:  reason,
		Message: message,
	}
}

// Unhealthy returns an unhealthy result for name.
//
// Unhealthy is the strongest negative health observation. It describes the
// checked scope only; it does not by itself prescribe restart, traffic removal,
// admission closure, or scheduler exclusion.
func Unhealthy(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusUnhealthy,
		Reason:  reason,
		Message: message,
	}
}

// Unknown returns an unknown result for name.
//
// Unknown means the checker could not produce a reliable health observation. It
// is useful for timeouts, cancellations, uninitialized cached results, missing
// state, invalid caller-controlled input, or inconclusive checks.
func Unknown(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusUnknown,
		Reason:  reason,
		Message: message,
	}
}

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
func (r Result) WithDuration(duration time.Duration) Result {
	r.Duration = duration
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
// Normalize is intended for checker, evaluator, registry, and report boundaries
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

// IsValid reports whether r is structurally valid as a health result.
//
// A valid result has a known Status value and a non-negative Duration. Name is
// not required for structural validity because the zero value is an unnamed
// StatusUnknown result and aggregators may fill the name from checker ownership.
func (r Result) IsValid() bool {
	return r.Status.IsValid() && r.Duration >= 0
}

// IsNamed reports whether r has a non-empty logical check name.
//
// Result names are used by registries, reports, tests, diagnostics, and adapters
// that expose individual check output. Aggregators SHOULD normalize unnamed
// checker results with the owning checker name.
func (r Result) IsNamed() bool {
	return r.Name != ""
}

// IsObserved reports whether r has an observation timestamp.
//
// A zero observation time means the result has not yet been timestamped by the
// checker or evaluator that owns the observation boundary.
func (r Result) IsObserved() bool {
	return !r.Observed.IsZero()
}

// HasCause reports whether r preserves an internal lower-level cause.
//
// Cause is useful for owner-controlled diagnostics and tests. Public adapters
// MUST NOT treat HasCause as permission to expose the cause.
func (r Result) HasCause() bool {
	return r.Cause != nil
}

// HasReason reports whether r has reason.
//
// HasReason performs exact reason matching. It intentionally does not interpret
// reason categories or status severity. Use Reason category helpers when callers
// need broader classification such as dependency, control, freshness, or
// observation reasons.
func (r Result) HasReason(reason Reason) bool {
	return r.Reason == reason
}

// IsAffirmative reports whether r is a positive health observation.
//
// This is a convenience wrapper over Status.IsAffirmative.
func (r Result) IsAffirmative() bool {
	return r.Status.IsAffirmative()
}

// IsNegative reports whether r is a strictly negative health observation.
//
// This is a convenience wrapper over Status.IsNegative.
func (r Result) IsNegative() bool {
	return r.Status.IsNegative()
}

// IsKnown reports whether r represents a determined operational state.
//
// This is a convenience wrapper over Status.IsKnown.
func (r Result) IsKnown() bool {
	return r.Status.IsKnown()
}

// IsOperational reports whether r describes a component that is still operating
// or making progress.
//
// This is a convenience wrapper over Status.IsOperational. Operational does not
// mean ready, admitted, routable, or fully healthy.
func (r Result) IsOperational() bool {
	return r.Status.IsOperational()
}

// MoreSevereThan reports whether r has a more severe status than other.
//
// MoreSevereThan compares only Status severity. It intentionally ignores Reason,
// Message, Duration, Observed, Cause, and Name because those fields do not define
// the core status ordering.
func (r Result) MoreSevereThan(other Result) bool {
	return r.Status.MoreSevereThan(other.Status)
}
