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

import (
	"testing"
	"time"
)

func TestReportValueSemantics(t *testing.T) {
	t.Parallel()

	report := Report{
		Target:   TargetReady,
		Status:   StatusDegraded,
		Observed: testObserved,
		Duration: time.Second,
		Checks: []Result{
			Healthy("storage"),
			Degraded("queue", ReasonOverloaded, "queue overloaded"),
		},
	}

	copied := report
	copied.Checks = append([]Result(nil), report.Checks...)
	copied.Checks[0] = Unhealthy("storage", ReasonFatal, "fatal")

	if report.Checks[0].Status != StatusHealthy {
		t.Fatalf("report check mutated through copy = %s, want healthy", report.Checks[0].Status)
	}
}
