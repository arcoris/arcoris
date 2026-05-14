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
	"testing"

	"arcoris.dev/health"
)

func TestAssertHelpersAcceptExpectedHealthValues(t *testing.T) {
	t.Parallel()

	report := MixedReport(health.TargetReady)
	AssertReportStatus(t, report, health.StatusUnhealthy)
	AssertReportTarget(t, report, health.TargetReady)
	AssertCheckOrder(t, report, "storage", "queue", "cache", "database")
	AssertReasons(t, report, health.ReasonOverloaded, health.ReasonNotObserved, health.ReasonDependencyUnavailable)
	AssertValidReport(t, report)

	invalid := report
	invalid.Duration = -1
	AssertInvalidReport(t, invalid)

	AssertResultStatus(t, report.Checks[0], health.StatusHealthy)
	AssertResultReason(t, report.Checks[1], health.ReasonOverloaded)
}
