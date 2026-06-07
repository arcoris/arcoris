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

import "testing"

func TestReportAccessorsUseReportOwnedChecks(t *testing.T) {
	t.Parallel()

	report := Report{
		Target:   TargetReady,
		Status:   StatusHealthy,
		Observed: testObserved,
		Checks:   []Result{Healthy("storage"), Healthy("queue")},
	}

	if !report.IsObserved() || report.Empty() || !report.Passed(ReadyPolicy()) {
		t.Fatalf("access predicates mismatch for %+v", report)
	}
	if res, ok := report.Check("queue"); !ok || res.Name != "queue" {
		t.Fatalf("Check(queue) = %+v, %v; want queue true", res, ok)
	}

	copy := report.ChecksCopy()
	copy[0] = Unhealthy("storage", ReasonFatal, "fatal")
	if report.Checks[0].Status != StatusHealthy {
		t.Fatal("ChecksCopy exposed mutable report storage")
	}
}
