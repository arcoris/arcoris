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

	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Check implements grpc.health.v1.Health.Check for one configured service.
//
// The method evaluates exactly one mapped package-health target through
// health.Evaluator and converts the resulting report status through the mapping
// policy. Unknown services and evaluator failures are returned as generic gRPC
// errors so raw health causes, panic details, credentials, or infrastructure
// addresses cannot leak through the transport boundary.
func (s *Server) Check(
	ctx context.Context,
	req *healthpb.HealthCheckRequest,
) (*healthpb.HealthCheckResponse, error) {
	if s == nil || nilSource(s.source) {
		return nil, status.Error(codes.Internal, healthServerUnavailableMessage)
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, nilCheckRequestMessage)
	}

	mapping, ok := s.service(req.GetService())
	if !ok {
		return nil, status.Error(codes.NotFound, unknownServiceMessage)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	report, err := s.source.Evaluate(ctx, mapping.Target)
	if err != nil {
		return nil, status.Error(codes.Internal, healthEvaluationFailedMessage)
	}

	return &healthpb.HealthCheckResponse{
		Status: ServingStatus(report.Status, mapping.Policy),
	}, nil
}
