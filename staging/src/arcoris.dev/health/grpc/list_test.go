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

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestListRejectsNilRequest(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))
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

	source := healthtest.NewTargetSource(map[health.Target]health.Report{
		health.TargetReady:   healthtest.HealthyReport(health.TargetReady),
		health.TargetStartup: healthtest.UnhealthyReport(health.TargetStartup),
		health.TargetLive:    healthtest.DegradedReport(health.TargetLive),
	})
	server := mustNewServer(t, source, WithTargetServices())

	resp, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if err != nil {
		t.Fatalf("List() = %v, want nil", err)
	}

	statuses := resp.GetStatuses()
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
		healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)),
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

	source := healthtest.NewCountingSource(map[health.Target]health.Report{
		health.TargetReady: healthtest.HealthyReport(health.TargetReady),
	})
	server := mustNewServer(t, source, WithService("ready-alt", health.TargetReady))

	resp, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if err != nil {
		t.Fatalf("List() = %v, want nil", err)
	}
	if len(resp.GetStatuses()) != 2 {
		t.Fatalf("statuses length = %d, want 2", len(resp.GetStatuses()))
	}
	if calls := source.Calls(health.TargetReady); calls != 1 {
		t.Fatalf("ready calls = %d, want 1", calls)
	}
}

func TestListEvaluationErrorMapsAffectedServicesToUnknown(t *testing.T) {
	t.Parallel()

	raw := errors.New("raw database outage")
	source := healthtest.NewCountingSource(map[health.Target]health.Report{
		health.TargetLive: healthtest.HealthyReport(health.TargetLive),
	})
	source.SetError(health.TargetReady, raw)
	server := mustNewServer(
		t,
		source,
		WithService("ready-alt", health.TargetReady),
		WithService("live", health.TargetLive),
	)

	resp, err := server.List(context.Background(), &healthpb.HealthListRequest{})
	if err != nil {
		t.Fatalf("List() = %v, want nil", err)
	}
	statuses := resp.GetStatuses()
	if statuses[""].GetStatus() != healthpb.HealthCheckResponse_UNKNOWN {
		t.Fatalf("default status = %s, want UNKNOWN", statuses[""].GetStatus())
	}
	if statuses["ready-alt"].GetStatus() != healthpb.HealthCheckResponse_UNKNOWN {
		t.Fatalf("ready-alt status = %s, want UNKNOWN", statuses["ready-alt"].GetStatus())
	}
	if statuses["live"].GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("live status = %s, want SERVING", statuses["live"].GetStatus())
	}
	if calls := source.Calls(health.TargetReady); calls != 1 {
		t.Fatalf("ready calls = %d, want 1", calls)
	}
}
