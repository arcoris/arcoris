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

// ChecksByReason returns check results with reason in report order.
//
// The returned slice is newly allocated. ReasonNone is matched like any other
// reason so callers can explicitly find results that did not provide a reason.
func (r Report) ChecksByReason(reason Reason) []Result {
	var matches []Result
	for _, res := range r.Checks {
		if res.HasReason(reason) {
			matches = append(matches, res)
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

	for _, res := range r.Checks {
		if res.Reason == ReasonNone {
			continue
		}
		if _, exists := seen[res.Reason]; exists {
			continue
		}

		seen[res.Reason] = struct{}{}
		reasons = append(reasons, res.Reason)
	}

	return reasons
}

// FailedChecks returns check results whose statuses fail policy.
func (r Report) FailedChecks(policy TargetPolicy) []Result {
	var failed []Result
	for _, res := range r.Checks {
		if policy.Fails(res.Status) {
			failed = append(failed, res)
		}
	}

	return failed
}

// DegradedChecks returns check results with StatusDegraded.
func (r Report) DegradedChecks() []Result {
	var degraded []Result
	for _, res := range r.Checks {
		if res.Status == StatusDegraded {
			degraded = append(degraded, res)
		}
	}

	return degraded
}

// UnknownChecks returns check results with StatusUnknown.
func (r Report) UnknownChecks() []Result {
	var unknown []Result
	for _, res := range r.Checks {
		if res.Status == StatusUnknown {
			unknown = append(unknown, res)
		}
	}

	return unknown
}
