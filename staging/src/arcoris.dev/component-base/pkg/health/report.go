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

// Report describes one target-level health evaluation.
//
// A Report is the Evaluator output for a concrete Target. It aggregates the most
// severe Status from the individual check Results while preserving every Result
// in deterministic registry order. Report intentionally remains a plain value so
// callers can store, copy, render, or adapt it without owning evaluator state.
//
// Report does not define transport behavior. HTTP status mapping, gRPC serving
// state, metrics, logging, restart policy, admission policy, routing, and
// scheduler decisions belong outside package health.
type Report struct {
	// Target is the health scope that was evaluated.
	Target Target

	// Status is the aggregate target status.
	//
	// Evaluator computes Status as the most severe Result status for the target.
	// A report with no checks uses StatusUnknown because no affirmative health
	// observation exists.
	Status Status

	// Observed is the time at which the report was produced.
	Observed time.Time

	// Duration is the evaluator-observed elapsed time for the target evaluation.
	//
	// Evaluator clamps negative durations to zero so mutable test clocks cannot
	// produce invalid runtime reports.
	Duration time.Duration

	// Checks contains normalized check Results in registry order.
	//
	// The slice is caller-owned after Evaluate returns. Evaluator does not retain
	// full reports or mutate returned Results.
	Checks []Result
}

// IsValid reports whether r is structurally valid as a health report.
//
// The zero report is valid and represents "not evaluated yet." Every non-zero or
// evaluated report must use a concrete Target. This prevents caller-constructed
// reports from carrying concrete check results under TargetUnknown.
func (r Report) IsValid() bool {
	if !r.Target.IsValid() || !r.Status.IsValid() || r.Duration < 0 {
		return false
	}

	if r.Target == TargetUnknown {
		return r.Status == StatusUnknown &&
			r.Observed.IsZero() &&
			r.Duration == 0 &&
			len(r.Checks) == 0
	}

	for _, check := range r.Checks {
		if !check.IsValid() {
			return false
		}
	}

	return true
}

// IsObserved reports whether r has a report-level observation timestamp.
func (r Report) IsObserved() bool {
	return !r.Observed.IsZero()
}

// Empty reports whether r contains no check results.
func (r Report) Empty() bool {
	return len(r.Checks) == 0
}

// Passed reports whether r.Status passes policy.
func (r Report) Passed(policy TargetPolicy) bool {
	return policy.Passes(r.Status)
}

// Failed reports whether r.Status fails policy.
func (r Report) Failed(policy TargetPolicy) bool {
	return policy.Fails(r.Status)
}

// Check returns the first check result with name.
//
// Reports preserve registry order and registry names are unique per target, so a
// report produced by Evaluator can contain at most one result for a valid check
// name. The method still returns the first match defensively because Report is a
// plain caller-owned value.
func (r Report) Check(name string) (Result, bool) {
	for _, check := range r.Checks {
		if check.Name == name {
			return check, true
		}
	}

	return Result{}, false
}

// HasReason reports whether at least one check result has reason.
//
// HasReason performs exact reason matching. It intentionally does not interpret
// reason categories or status severity. Use Reason category helpers on individual
// results when callers need broader classification.
func (r Report) HasReason(reason Reason) bool {
	for _, check := range r.Checks {
		if check.HasReason(reason) {
			return true
		}
	}

	return false
}

// ChecksCopy returns a defensive copy of r.Checks.
func (r Report) ChecksCopy() []Result {
	if len(r.Checks) == 0 {
		return nil
	}

	checks := make([]Result, len(r.Checks))
	copy(checks, r.Checks)

	return checks
}

// ChecksByReason returns check results with reason in report order.
//
// The returned slice is newly allocated. ReasonNone is matched like any other
// reason so callers can explicitly find results that did not provide a reason.
func (r Report) ChecksByReason(reason Reason) []Result {
	var matches []Result
	for _, check := range r.Checks {
		if check.HasReason(reason) {
			matches = append(matches, check)
		}
	}

	return matches
}

// Reasons returns unique non-empty reasons in first-seen report order.
//
// ReasonNone is omitted because it represents the absence of a reason rather
// than a diagnostic classification. Use ChecksByReason(ReasonNone) when callers
// need to inspect results that intentionally provided no reason.
func (r Report) Reasons() []Reason {
	var reasons []Reason
	seen := make(map[Reason]struct{})

	for _, check := range r.Checks {
		if check.Reason == ReasonNone {
			continue
		}
		if _, exists := seen[check.Reason]; exists {
			continue
		}

		seen[check.Reason] = struct{}{}
		reasons = append(reasons, check.Reason)
	}

	return reasons
}

// FailedChecks returns check results whose statuses fail policy.
func (r Report) FailedChecks(policy TargetPolicy) []Result {
	var failed []Result
	for _, check := range r.Checks {
		if policy.Fails(check.Status) {
			failed = append(failed, check)
		}
	}

	return failed
}

// DegradedChecks returns check results with StatusDegraded.
func (r Report) DegradedChecks() []Result {
	var degraded []Result
	for _, check := range r.Checks {
		if check.Status == StatusDegraded {
			degraded = append(degraded, check)
		}
	}

	return degraded
}

// UnknownChecks returns check results with StatusUnknown.
func (r Report) UnknownChecks() []Result {
	var unknown []Result
	for _, check := range r.Checks {
		if check.Status == StatusUnknown {
			unknown = append(unknown, check)
		}
	}

	return unknown
}
