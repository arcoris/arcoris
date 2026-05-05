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
	"testing"

	"arcoris.dev/component-base/pkg/health"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestListRejectsNilRequest(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, staticSource{status: health.StatusHealthy})
	_, err := server.List(context.Background(), nil)
	if grpcCode(err) != codes.InvalidArgument {
		t.Fatalf("List(nil) code = %s, want InvalidArgument", grpcCode(err))
	}
}

func TestListRejectsNilServer(t *testing.T) {
	t.Parallel()

	var server *Server
	_, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if grpcCode(err) != codes.Internal {
		t.Fatalf("List(nil server) code = %s, want Internal", grpcCode(err))
	}
}

func TestListReturnsConfiguredServices(t *testing.T) {
	t.Parallel()

	source := targetSource{statuses: map[health.Target]health.Status{
		health.TargetReady:   health.StatusHealthy,
		health.TargetStartup: health.StatusUnhealthy,
		health.TargetLive:    health.StatusDegraded,
	}}
	server := mustNewServer(t, source, WithTargetServices())

	response, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if err != nil {
		t.Fatalf("List() = %v, want nil", err)
	}

	statuses := response.GetStatuses()
	if len(statuses) != 4 {
		t.Fatalf("statuses length = %d, want 4", len(statuses))
	}
	if statuses[""].GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("default status = %s, want SERVING", statuses[""].GetStatus())
	}
	if statuses["startup"].GetStatus() != healthpb.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("startup status = %s, want NOT_SERVING", statuses["startup"].GetStatus())
	}
	if statuses["live"].GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("live status = %s, want SERVING", statuses["live"].GetStatus())
	}
	if statuses["ready"].GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("ready status = %s, want SERVING", statuses["ready"].GetStatus())
	}
}

func TestListMaxServiceLimit(t *testing.T) {
	t.Parallel()

	server := mustNewServer(
		t,
		staticSource{status: health.StatusHealthy},
		WithTargetServices(),
		WithMaxListServices(1),
	)

	_, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if grpcCode(err) != codes.ResourceExhausted {
		t.Fatalf("List() code = %s, want ResourceExhausted", grpcCode(err))
	}
}

func TestListEvaluatesSharedTargetOnce(t *testing.T) {
	t.Parallel()

	source := newCountingSource()
	source.statuses[health.TargetReady] = health.StatusHealthy
	server := mustNewServer(t, source, WithService("ready-alt", health.TargetReady))

	response, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if err != nil {
		t.Fatalf("List() = %v, want nil", err)
	}
	if len(response.GetStatuses()) != 2 {
		t.Fatalf("statuses length = %d, want 2", len(response.GetStatuses()))
	}
	if calls := source.callsFor(health.TargetReady); calls != 1 {
		t.Fatalf("ready calls = %d, want 1", calls)
	}
}

func TestListEvaluationErrorMapsAffectedServicesToUnknown(t *testing.T) {
	t.Parallel()

	raw := errors.New("raw database outage")
	source := newCountingSource()
	source.errors[health.TargetReady] = raw
	source.statuses[health.TargetLive] = health.StatusHealthy
	server := mustNewServer(
		t,
		source,
		WithService("ready-alt", health.TargetReady),
		WithService("live", health.TargetLive),
	)

	response, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if err != nil {
		t.Fatalf("List() = %v, want nil", err)
	}
	statuses := response.GetStatuses()
	if statuses[""].GetStatus() != healthpb.HealthCheckResponse_UNKNOWN {
		t.Fatalf("default status = %s, want UNKNOWN", statuses[""].GetStatus())
	}
	if statuses["ready-alt"].GetStatus() != healthpb.HealthCheckResponse_UNKNOWN {
		t.Fatalf("ready-alt status = %s, want UNKNOWN", statuses["ready-alt"].GetStatus())
	}
	if statuses["live"].GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("live status = %s, want SERVING", statuses["live"].GetStatus())
	}
	if calls := source.callsFor(health.TargetReady); calls != 1 {
		t.Fatalf("ready calls = %d, want 1", calls)
	}
}
