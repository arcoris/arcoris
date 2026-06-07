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

// IsNamed reports whether r has a non-empty logical check name.
//
// Result names are used by resolvers, reports, tests, diagnostics, and adapters
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
