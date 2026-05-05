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

import healthpb "google.golang.org/grpc/health/grpc_health_v1"

// Server implements the standard grpc.health.v1.Health service over Source.
//
// A Server is intentionally a thin adapter. Its immutable configuration maps
// transport service names to package-health targets and policies. It does not
// cache reports, maintain serving state, or mutate the Source; every RPC that
// needs health state asks Source to evaluate the configured target.
type Server struct {
	// UnimplementedHealthServer keeps the generated service forward-compatible
	// if the standard gRPC health API grows additional methods.
	healthpb.UnimplementedHealthServer

	// source is the only runtime health boundary. The adapter never reaches into
	// registries or check lists directly.
	source Source

	// services provides immutable lookup by exact gRPC service name.
	services map[string]ServiceMapping

	// order preserves caller configuration order for Services and List.
	order []string

	// config keeps validated adapter-local settings needed after construction,
	// such as Watch cadence, List guardrails, and clock ownership.
	config config
}
