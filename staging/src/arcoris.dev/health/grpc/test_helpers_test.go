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
	"io"
	"sync"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// testTimeout bounds asynchronous Watch test waits.
const testTimeout = 5 * time.Second

// watchStream is a minimal Health_WatchServer test double.
type watchStream struct {
	// ctx is returned by Context and canceled by tests to end Watch.
	ctx context.Context

	// cancel stops ctx.
	cancel context.CancelFunc

	// mu protects sendErr and responses.
	mu sync.Mutex

	// sendErr forces Send to fail when non-nil.
	sendErr error

	// responses stores sent responses for postconditions.
	responses []*healthpb.HealthCheckResponse

	// sent streams responses to tests without polling responses under mu.
	sent chan *healthpb.HealthCheckResponse
}

// newWatchStream returns an initialized Watch stream test double.
func newWatchStream() *watchStream {
	ctx, cancel := context.WithCancel(context.Background())
	return &watchStream{
		ctx:    ctx,
		cancel: cancel,
		sent:   make(chan *healthpb.HealthCheckResponse, 16),
	}
}

// Send records a detached response or returns the configured send error.
func (s *watchStream) Send(response *healthpb.HealthCheckResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sendErr != nil {
		return s.sendErr
	}

	copied := &healthpb.HealthCheckResponse{Status: response.GetStatus()}
	s.responses = append(s.responses, copied)
	select {
	case s.sent <- copied:
	default:
	}

	return nil
}

// SetHeader satisfies grpc.ServerStream.
func (s *watchStream) SetHeader(metadata.MD) error {
	return nil
}

// SendHeader satisfies grpc.ServerStream.
func (s *watchStream) SendHeader(metadata.MD) error {
	return nil
}

// SetTrailer satisfies grpc.ServerStream.
func (s *watchStream) SetTrailer(metadata.MD) {}

// Context returns the stream context observed by Watch.
func (s *watchStream) Context() context.Context {
	return s.ctx
}

// SendMsg satisfies grpc.ServerStream for the generated stream interface.
func (s *watchStream) SendMsg(any) error {
	return nil
}

// RecvMsg satisfies grpc.ServerStream for the generated stream interface.
func (s *watchStream) RecvMsg(any) error {
	return io.EOF
}

// mustNewServer builds a Server or fails the test.
func mustNewServer(t *testing.T, source Source, options ...Option) *Server {
	t.Helper()

	server, err := NewServer(source, options...)
	if err != nil {
		t.Fatalf("NewServer() = %v, want nil", err)
	}

	return server
}

// grpcCode returns the canonical gRPC code for err.
func grpcCode(err error) codes.Code {
	return status.Code(err)
}

// testClock returns a FakeClock anchored at the shared health fixture time.
func testClock() *clock.FakeClock {
	return clock.NewFakeClock(healthtest.ObservedTime)
}

// mustReceiveStatus waits for the next Watch response and returns its status.
func mustReceiveStatus(
	t *testing.T,
	stream *watchStream,
) healthpb.HealthCheckResponse_ServingStatus {
	t.Helper()

	select {
	case response := <-stream.sent:
		return response.GetStatus()
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for watch response")
		return healthpb.HealthCheckResponse_UNKNOWN
	}
}

// assertNoWatchStatus verifies that Watch did not send another response.
func assertNoWatchStatus(t *testing.T, stream *watchStream) {
	t.Helper()

	select {
	case response := <-stream.sent:
		t.Fatalf("unexpected watch response: %s", response.GetStatus())
	default:
	}
}

// waitForWatchDone waits for a Watch goroutine to return.
func waitForWatchDone(t *testing.T, done <-chan error) error {
	t.Helper()

	select {
	case err := <-done:
		return err
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for Watch to return")
		return nil
	}
}

// sourceCallCounter is the healthtest source call-count interface used by Watch tests.
type sourceCallCounter interface {
	Calls(health.Target) int
}

// waitForSourceCalls waits until source has observed at least want calls.
func waitForSourceCalls(t *testing.T, source sourceCallCounter, target health.Target, want int) {
	t.Helper()

	if source.Calls(target) >= want {
		return
	}

	deadline := time.After(testTimeout)
	for {
		if source.Calls(target) >= want {
			return
		}
		select {
		case <-time.After(time.Millisecond):
		case <-deadline:
			t.Fatalf("timed out waiting for source calls target=%s >= %d", target, want)
		}
	}
}
