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
	"errors"
	"strings"
	"testing"

	"arcoris.dev/component-base/pkg/health"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestCheckRejectsNilRequest(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, staticSource{status: health.StatusHealthy})
	_, err := server.Check(context.Background(), nil)
	if grpcCode(err) != codes.InvalidArgument {
		t.Fatalf("Check(nil) code = %s, want InvalidArgument", grpcCode(err))
	}
}

func TestCheckRejectsNilServer(t *testing.T) {
	t.Parallel()

	var server *Server
	_, err := server.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if grpcCode(err) != codes.Internal {
		t.Fatalf("Check(nil server) code = %s, want Internal", grpcCode(err))
	}
}

func TestCheckRejectsUnknownService(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, staticSource{status: health.StatusHealthy})
	_, err := server.Check(context.Background(), &healthpb.HealthCheckRequest{Service: "missing"})
	if grpcCode(err) != codes.NotFound {
		t.Fatalf("Check(missing) code = %s, want NotFound", grpcCode(err))
	}
	if status.Convert(err).Message() != "unknown service" {
		t.Fatalf("message = %q, want unknown service", status.Convert(err).Message())
	}
}

func TestCheckMapsReports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		source  Source
		options []Option
		service string
		want    healthpb.HealthCheckResponse_ServingStatus
	}{
		{
			name:   "healthy ready",
			source: staticSource{status: health.StatusHealthy},
			want:   healthpb.HealthCheckResponse_SERVING,
		},
		{
			name:   "unhealthy ready",
			source: staticSource{status: health.StatusUnhealthy},
			want:   healthpb.HealthCheckResponse_NOT_SERVING,
		},
		{
			name:   "unknown report",
			source: staticSource{status: health.StatusUnknown},
			want:   healthpb.HealthCheckResponse_UNKNOWN,
		},
		{
			name:   "degraded ready",
			source: staticSource{status: health.StatusDegraded},
			want:   healthpb.HealthCheckResponse_NOT_SERVING,
		},
		{
			name:    "degraded live",
			source:  targetSource{statuses: map[health.Target]health.Status{health.TargetLive: health.StatusDegraded}},
			options: []Option{WithTargetServices()},
			service: "live",
			want:    healthpb.HealthCheckResponse_SERVING,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := mustNewServer(t, tc.source, tc.options...)
			response, err := server.Check(context.Background(), &healthpb.HealthCheckRequest{Service: tc.service})
			if err != nil {
				t.Fatalf("Check() = %v, want nil", err)
			}
			if response.GetStatus() != tc.want {
				t.Fatalf("Status = %s, want %s", response.GetStatus(), tc.want)
			}
		})
	}
}

func TestCheckSourceErrorIsGeneric(t *testing.T) {
	t.Parallel()

	raw := errors.New("database password=secret is down")
	server := mustNewServer(t, errorSource{err: raw})

	_, err := server.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if grpcCode(err) != codes.Internal {
		t.Fatalf("Check() code = %s, want Internal", grpcCode(err))
	}
	message := status.Convert(err).Message()
	if message != "health evaluation failed" {
		t.Fatalf("message = %q, want health evaluation failed", message)
	}
	if strings.Contains(message, "password=secret") {
		t.Fatalf("message leaked raw source error: %q", message)
	}
}
