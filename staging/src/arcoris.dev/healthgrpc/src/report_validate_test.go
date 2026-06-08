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

package healthgrpc

import (
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
)

func TestValidReportForTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		report health.Report
		want   bool
	}{
		{
			name:   "valid",
			report: healthtest.HealthyReport(health.TargetReady),
			want:   true,
		},
		{
			name:   "wrong target",
			report: healthtest.HealthyReport(health.TargetLive),
		},
		{
			name:   "invalid",
			report: grpcInvalidReport(health.TargetReady),
		},
		{
			name:   "inconsistent",
			report: grpcInconsistentReport(health.TargetReady),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := validReportForTarget(tc.report, health.TargetReady); got != tc.want {
				t.Fatalf("validReportForTarget() = %v, want %v", got, tc.want)
			}
		})
	}
}

func grpcInvalidReport(target health.Target) health.Report {
	report := healthtest.HealthyReport(target)
	report.Status = health.Status(255)

	return report
}

func grpcInconsistentReport(target health.Target) health.Report {
	return healthtest.Report(
		target,
		health.StatusHealthy,
		health.Unhealthy("database", health.ReasonFatal, "database failed"),
	)
}
