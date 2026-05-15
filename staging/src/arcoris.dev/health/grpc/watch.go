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

// Watch implements grpc.health.v1.Health.Watch as per-stream polling.
//
// Each call owns its ticker and stream lifetime. The adapter sends the initial
// status immediately, then sends only status changes. Unknown services follow
// the standard gRPC health behavior by sending SERVICE_UNKNOWN once and then
// waiting for the stream to end.
func (s *Server) Watch(
	req *healthpb.HealthCheckRequest,
	stream grpc.ServerStreamingServer[healthpb.HealthCheckResponse],
) error {
	if s == nil || nilSource(s.source) || stream == nil {
		return status.Error(codes.Canceled, watchEndedMessage)
	}
	if req == nil {
		return status.Error(codes.InvalidArgument, nilWatchRequestMessage)
	}

	mapping, ok := s.service(req.GetService())
	if !ok {
		if err := sendWatchStatus(stream, healthpb.HealthCheckResponse_SERVICE_UNKNOWN); err != nil {
			return err
		}
		return waitForWatchEnd(stream.Context())
	}

	ctx := stream.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	last := s.watchServingStatus(ctx, mapping)
	if err := sendWatchStatus(stream, last); err != nil {
		return err
	}

	ticker := s.config.clock.NewTicker(s.config.watchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return status.Error(codes.Canceled, watchEndedMessage)
		case <-ticker.C():
			next := s.watchServingStatus(ctx, mapping)
			if next == last {
				continue
			}
			if err := sendWatchStatus(stream, next); err != nil {
				return err
			}
			last = next
		}
	}
}
