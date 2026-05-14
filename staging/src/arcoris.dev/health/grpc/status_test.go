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

package healthgrpc

import (
	"testing"

	"arcoris.dev/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestServingStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status health.Status
		policy health.TargetPolicy
		want   healthpb.HealthCheckResponse_ServingStatus
	}{
		{"healthy", health.StatusHealthy, health.ReadyPolicy(), healthpb.HealthCheckResponse_SERVING},
		{"starting startup", health.StatusStarting, health.StartupPolicy(), healthpb.HealthCheckResponse_NOT_SERVING},
		{"starting ready", health.StatusStarting, health.ReadyPolicy(), healthpb.HealthCheckResponse_NOT_SERVING},
		{"starting live", health.StatusStarting, health.LivePolicy(), healthpb.HealthCheckResponse_SERVING},
		{"degraded startup", health.StatusDegraded, health.StartupPolicy(), healthpb.HealthCheckResponse_NOT_SERVING},
		{"degraded ready", health.StatusDegraded, health.ReadyPolicy(), healthpb.HealthCheckResponse_NOT_SERVING},
		{"degraded live", health.StatusDegraded, health.LivePolicy(), healthpb.HealthCheckResponse_SERVING},
		{
			name:   "degraded custom allow",
			status: health.StatusDegraded,
			policy: health.ReadyPolicy().WithDegraded(true),
			want:   healthpb.HealthCheckResponse_SERVING,
		},
		{"unknown", health.StatusUnknown, health.LivePolicy(), healthpb.HealthCheckResponse_UNKNOWN},
		{"unhealthy", health.StatusUnhealthy, health.LivePolicy(), healthpb.HealthCheckResponse_NOT_SERVING},
		{"invalid", health.Status(255), health.LivePolicy(), healthpb.HealthCheckResponse_UNKNOWN},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := ServingStatus(tc.status, tc.policy); got != tc.want {
				t.Fatalf("ServingStatus() = %s, want %s", got, tc.want)
			}
		})
	}
}
