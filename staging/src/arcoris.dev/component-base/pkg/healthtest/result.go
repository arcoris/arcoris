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
	"errors"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

// ObservedTime is the stable observation timestamp used by healthtest fixtures.
//
// Tests in adapters often compare rendered timestamps, response DTOs, cached
// snapshots, and report validity. A single exported timestamp keeps those tests
// deterministic without each package inventing its own local clock constant.
var ObservedTime = time.Date(2026, 5, 4, 12, 30, 15, 123456789, time.UTC)

// HealthyResult returns a deterministic healthy result for name.
//
// The returned value is already observed at ObservedTime. It is suitable for
// evaluator, report, adapter-rendering, and cache tests that need a complete
// positive check result without constructing one field by field.
func HealthyResult(name string) health.Result {
	return ResultWithObservation(health.Healthy(name))
}

// StartingResult returns a deterministic starting result for name.
//
// ReasonStarting is used deliberately so tests can assert both StatusStarting
// and reason propagation through reports and adapters.
func StartingResult(name string) health.Result {
	return ResultWithObservation(health.Starting(name, health.ReasonStarting, "health check is starting"))
}

// DegradedResult returns a deterministic degraded result for name and reason.
//
// The caller supplies the reason because degraded results are commonly used to
// exercise target-policy differences such as readiness rejecting degradation
// while liveness accepts it.
func DegradedResult(name string, reason health.Reason) health.Result {
	return ResultWithObservation(health.Degraded(name, reason, "health check is degraded"))
}

// UnhealthyResult returns a deterministic unhealthy result for name and reason.
//
// The message is intentionally safe and generic. Tests that need private
// diagnostics should attach a cause with ResultWithPrivateCause instead of
// hiding sensitive text in Message.
func UnhealthyResult(name string, reason health.Reason) health.Result {
	return ResultWithObservation(health.Unhealthy(name, reason, "health check is unhealthy"))
}

// UnknownResult returns a deterministic unknown result for name and reason.
//
// Unknown fixtures are useful for not-yet-observed, canceled, timeout, and
// source-error paths where adapters must preserve UNKNOWN instead of collapsing
// it into a negative serving state.
func UnknownResult(name string, reason health.Reason) health.Result {
	return ResultWithObservation(health.Unknown(name, reason, "health check is unknown"))
}

// ResultWithObservation returns result with the canonical observation time.
//
// The helper keeps the rest of the result untouched so tests can build unusual
// or intentionally invalid values while still avoiding zero timestamps when the
// observation time is not the behavior under test.
func ResultWithObservation(result health.Result) health.Result {
	return result.WithObserved(ObservedTime)
}

// ResultWithDuration returns result with duration set.
//
// Duration is often rendered by transport adapters, so this helper makes those
// fixtures explicit without requiring each test to remember the fluent Result
// method name.
func ResultWithDuration(result health.Result, duration time.Duration) health.Result {
	return result.WithDuration(duration)
}

// ResultWithPrivateCause returns result with an internal cause for leakage tests.
//
// The cause text is intentionally recognizable. Public adapters should never
// expose Result.Cause, so tests can scan rendered output for "private cause" and
// fail loudly if an adapter accidentally leaks internal diagnostics.
func ResultWithPrivateCause(result health.Result) health.Result {
	return result.WithCause(errors.New("private cause"))
}
