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

import "testing"

func TestAggregateStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		results []Result
		want    Status
	}{
		{name: "empty", want: StatusUnknown},
		{name: "healthy", results: []Result{Healthy("healthy")}, want: StatusHealthy},
		{name: "degraded", results: []Result{Healthy("healthy"), Degraded("degraded", ReasonOverloaded, "degraded")}, want: StatusDegraded},
		{name: "unknown", results: []Result{Degraded("degraded", ReasonOverloaded, "degraded"), Unknown("unknown", ReasonNotObserved, "unknown")}, want: StatusUnknown},
		{name: "unhealthy", results: []Result{Unknown("unknown", ReasonNotObserved, "unknown"), Unhealthy("unhealthy", ReasonFatal, "unhealthy")}, want: StatusUnhealthy},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := aggregateStatus(test.results); got != test.want {
				t.Fatalf("aggregateStatus() = %s, want %s", got, test.want)
			}
		})
	}
}
