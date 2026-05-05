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

package healthtest

import (
	"time"

	"arcoris.dev/component-base/pkg/health"
)

// Report returns a deterministic report for target, status, and checks.
//
// The returned report owns its Checks slice. Check results with zero observation
// time are stamped with ObservedTime so fixtures remain stable across packages.
// The helper does not recompute status from checks; callers pass the aggregate
// status explicitly so tests can construct edge cases around aggregation,
// adapter mapping, and validation behavior.
func Report(target health.Target, status health.Status, checks ...health.Result) health.Report {
	return health.Report{
		Target:   target,
		Status:   status,
		Observed: ObservedTime,
		Duration: 25 * time.Millisecond,
		Checks:   observedChecksCopy(checks),
	}
}

// HealthyReport returns a valid healthy report for target.
//
// The check name is derived from target so packages that assert names can use
// the fixture without adding their own naming convention.
func HealthyReport(target health.Target) health.Report {
	return Report(target, health.StatusHealthy, HealthyResult(target.String()+"_check"))
}

// StartingReport returns a valid starting report for target.
//
// The report is structurally valid and observed, making it suitable for adapter
// status-mapping tests rather than only value-constructor tests.
func StartingReport(target health.Target) health.Report {
	return Report(target, health.StatusStarting, StartingResult(target.String()+"_check"))
}

// DegradedReport returns a valid degraded report for target.
//
// ReasonOverloaded is used as a representative policy-relevant degradation
// reason. Tests that need a different reason can call Report directly.
func DegradedReport(target health.Target) health.Report {
	return Report(target, health.StatusDegraded, DegradedResult(target.String()+"_check", health.ReasonOverloaded))
}

// UnhealthyReport returns a valid unhealthy report for target.
//
// ReasonFatal keeps the fixture clearly negative without implying dependency,
// transport, or lifecycle behavior owned by another package.
func UnhealthyReport(target health.Target) health.Report {
	return Report(target, health.StatusUnhealthy, UnhealthyResult(target.String()+"_check", health.ReasonFatal))
}

// UnknownReport returns a valid unknown report for target.
//
// Unknown reports are used for cache misses, source failures, and not-yet-known
// health states. They remain observed so tests can distinguish "unknown status"
// from "zero report" for concrete targets. TargetUnknown returns the valid zero
// unknown report shape defined by package health.
func UnknownReport(target health.Target) health.Report {
	if target == health.TargetUnknown {
		return health.Report{Target: health.TargetUnknown, Status: health.StatusUnknown}
	}

	return Report(target, health.StatusUnknown, UnknownResult(target.String()+"_check", health.ReasonNotObserved))
}

// MixedReport returns a deterministic report with mixed check statuses.
//
// The fixture is useful for adapter rendering tests. It includes one private
// cause on the database result so public response tests can verify that
// transport adapters do not leak Result.Cause.
func MixedReport(target health.Target) health.Report {
	return Report(
		target,
		health.StatusUnhealthy,
		ResultWithDuration(HealthyResult("storage"), 2*time.Millisecond),
		ResultWithDuration(
			health.Degraded("queue", health.ReasonOverloaded, "queue is above soft capacity").
				WithObserved(ObservedTime),
			3*time.Millisecond,
		),
		ResultWithDuration(
			health.Unknown("cache", health.ReasonNotObserved, "cache has not reported yet").
				WithObserved(ObservedTime),
			4*time.Millisecond,
		),
		ResultWithPrivateCause(ResultWithDuration(
			health.Unhealthy(
				"database",
				health.ReasonDependencyUnavailable,
				"database dependency is unavailable",
			).WithObserved(ObservedTime),
			5*time.Millisecond,
		)),
	)
}

// observedChecksCopy returns a detached copy of checks with stable observations.
//
// Report fixtures are plain values that tests may mutate. Copying here prevents
// mutations of input slices from changing the returned report and avoids hidden
// aliasing between callers.
func observedChecksCopy(checks []health.Result) []health.Result {
	if len(checks) == 0 {
		return nil
	}

	copied := make([]health.Result, len(checks))
	copy(copied, checks)
	for i := range copied {
		if copied[i].Observed.IsZero() {
			copied[i] = copied[i].WithObserved(ObservedTime)
		}
	}

	return copied
}
