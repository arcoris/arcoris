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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// watchServingStatus evaluates mapping and returns the gRPC Watch status.
//
// Evaluator failures are reported as UNKNOWN instead of stream errors. This keeps
// Watch alive across transient evaluation failures and avoids exposing raw
// errors to clients.
func (s *Server) watchServingStatus(
	ctx context.Context,
	mapping ServiceMapping,
) healthpb.HealthCheckResponse_ServingStatus {
	report, err := s.source.Evaluate(ctx, mapping.Target)
	if err != nil {
		return healthpb.HealthCheckResponse_UNKNOWN
	}

	return ServingStatus(report.Status, mapping.Policy)
}

// sendWatchStatus sends one Watch response with generic error normalization.
//
// A stream Send failure can contain transport-specific details. The adapter
// returns a stable CANCELED status instead of forwarding those details.
func sendWatchStatus(
	stream grpc.ServerStreamingServer[healthpb.HealthCheckResponse],
	servingStatus healthpb.HealthCheckResponse_ServingStatus,
) error {
	if err := stream.Send(&healthpb.HealthCheckResponse{Status: servingStatus}); err != nil {
		return status.Error(codes.Canceled, watchEndedMessage)
	}

	return nil
}

// waitForWatchEnd blocks until an unknown-service Watch stream is canceled.
//
// The gRPC health protocol expects SERVICE_UNKNOWN to be delivered as a stream
// response rather than as a terminal NotFound error. After sending it, the
// adapter owns no further work and waits only for the caller's stream context.
func waitForWatchEnd(ctx context.Context) error {
	if ctx != nil {
		<-ctx.Done()
	}

	return status.Error(codes.Canceled, watchEndedMessage)
}
