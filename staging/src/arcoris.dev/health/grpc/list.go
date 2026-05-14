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
	"context"

	"arcoris.dev/health"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// List implements grpc.health.v1.Health.List for configured services.
//
// The method evaluates each unique package-health target at most once and then
// projects the target result onto every configured service that uses it. The
// protobuf response is a map, but the adapter still walks the immutable service
// order so target evaluation and tests remain deterministic.
func (s *Server) List(
	ctx context.Context,
	request *healthpb.HealthListRequest,
) (*healthpb.HealthListResponse, error) {
	if s == nil || nilSource(s.source) {
		return nil, status.Error(codes.Internal, healthServerUnavailableMessage)
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, nilListRequestMessage)
	}
	if len(s.order) > s.config.maxListServices {
		return nil, status.Error(codes.ResourceExhausted, tooManyHealthServicesMessage)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	byTarget := make(map[health.Target]targetEvaluation, len(s.order))
	statuses := make(map[string]*healthpb.HealthCheckResponse, len(s.order))

	for _, service := range s.order {
		mapping := s.services[service]
		result, ok := byTarget[mapping.Target]
		if !ok {
			result = s.evaluateTarget(ctx, mapping.Target)
			byTarget[mapping.Target] = result
		}

		servingStatus := healthpb.HealthCheckResponse_UNKNOWN
		if !result.failed {
			servingStatus = ServingStatus(result.status, mapping.Policy)
		}
		statuses[service] = &healthpb.HealthCheckResponse{Status: servingStatus}
	}

	return &healthpb.HealthListResponse{Statuses: statuses}, nil
}
