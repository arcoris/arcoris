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

package eval

import (
	"errors"
	"time"

	"arcoris.dev/health"
)

// invalidCheckerResult converts a malformed checker contract into a conservative
// result without executing Check.
func invalidCheckerResult(name string, err error, observed time.Time, d time.Duration) health.Result {
	resultName := ""
	reason := health.ReasonMisconfigured
	message := "health checker has invalid name"
	if errors.Is(err, health.ErrNilChecker) {
		reason = health.ReasonNotObserved
		message = "health checker is nil"
	}

	if health.ValidCheckName(name) {
		resultName = name
	}

	return health.Unknown(resultName, reason, message).
		WithCause(err).
		WithObserved(observed).
		WithDuration(nonNegativeDuration(d))
}

// normalizeEvaluatedResult applies evaluator-owned boundary normalization.
//
// Evaluator owns checker identity in reports. A checker may leave health.Result.Name
// empty, but it must not return another checker's name. Invalid reasons are also
// converted into unknown misconfiguration results so Evaluator never returns an
// invalid health.Report because of malformed checker output.
func normalizeEvaluatedResult(res health.Result, defaultName string, observed time.Time, d time.Duration) health.Result {
	d = nonNegativeDuration(d)

	if res.Name != "" && res.Name != defaultName {
		return health.Unknown(
			defaultName,
			health.ReasonMisconfigured,
			"health check returned a mismatched result name",
		).WithCause(MismatchedCheckResultError{
			CheckName:  defaultName,
			ResultName: res.Name,
		}).WithObserved(observed).WithDuration(d)
	}

	if !res.Reason.IsValid() {
		return health.Unknown(
			defaultName,
			health.ReasonMisconfigured,
			"health check returned an invalid reason",
		).WithCause(InvalidCheckResultError{
			CheckName: defaultName,
			Result:    res,
		}).WithObserved(observed).WithDuration(d)
	}

	res = res.Normalize(defaultName, observed)

	if res.Duration == 0 {
		res.Duration = d
	}

	if res.Duration < 0 {
		res.Duration = 0
	}

	return res
}

// nonNegativeDuration returns duration unless it is negative.
//
// Negative durations can occur with mutable fake clocks. Runtime reports should
// remain conservative and never expose negative elapsed time.
func nonNegativeDuration(d time.Duration) time.Duration {
	if d < 0 {
		return 0
	}

	return d
}
