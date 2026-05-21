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
	"arcoris.dev/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// ServingStatus maps a package-health status to a gRPC health serving status.
//
// StatusUnknown maps to gRPC UNKNOWN even though it fails every
// health.TargetPolicy. This preserves the difference between "not serving" and
// "not determined" at the transport boundary.
func ServingStatus(
	status health.Status,
	policy health.TargetPolicy,
) healthpb.HealthCheckResponse_ServingStatus {
	switch status {
	case health.StatusHealthy:
		return healthpb.HealthCheckResponse_SERVING
	case health.StatusStarting:
		if policy.AllowStarting {
			return healthpb.HealthCheckResponse_SERVING
		}
		return healthpb.HealthCheckResponse_NOT_SERVING
	case health.StatusDegraded:
		if policy.AllowDegraded {
			return healthpb.HealthCheckResponse_SERVING
		}
		return healthpb.HealthCheckResponse_NOT_SERVING
	case health.StatusUnhealthy:
		return healthpb.HealthCheckResponse_NOT_SERVING
	case health.StatusUnknown:
		return healthpb.HealthCheckResponse_UNKNOWN
	default:
		return healthpb.HealthCheckResponse_UNKNOWN
	}
}
