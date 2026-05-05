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
	"io"
	"sync"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/clock"
	"arcoris.dev/component-base/pkg/health"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// testTimeout bounds asynchronous Watch test waits.
const testTimeout = 5 * time.Second

// testObserved is the stable timestamp used by synthetic health reports.
var testObserved = time.Unix(100, 0)

// sourceFunc adapts a function to Source for focused tests.
type sourceFunc func(context.Context, health.Target) (health.Report, error)

// Evaluate calls fn with the supplied context and target.
func (fn sourceFunc) Evaluate(ctx context.Context, target health.Target) (health.Report, error) {
	return fn(ctx, target)
}

// staticSource returns one status for every target.
type staticSource struct {
	// status is the report status returned by Evaluate.
	status health.Status
}

// Evaluate returns a ready-target report with the configured status.
func (s staticSource) Evaluate(context.Context, health.Target) (health.Report, error) {
	return health.Report{Target: health.TargetReady, Status: s.status, Observed: testObserved}, nil
}

// targetSource returns statuses keyed by package-health target.
type targetSource struct {
	// statuses maps each target to the status returned for that target.
	statuses map[health.Target]health.Status
}

// Evaluate returns the target-specific status or UNKNOWN when absent.
func (s targetSource) Evaluate(_ context.Context, target health.Target) (health.Report, error) {
	status, ok := s.statuses[target]
	if !ok {
		status = health.StatusUnknown
	}

	return testReport(target, status), nil
}

// errorSource returns a stable evaluation error.
type errorSource struct {
	// err is returned by Evaluate; a default error is created when nil.
	err error
}

// Evaluate returns err without exposing a report.
func (s errorSource) Evaluate(context.Context, health.Target) (health.Report, error) {
	if s.err == nil {
		s.err = errors.New("source error")
	}

	return health.Report{}, s.err
}

// scriptedResult is one scripted Source response.
type scriptedResult struct {
	// status is returned as a report status when err is nil.
	status health.Status

	// err is returned instead of a report when non-nil.
	err error
}

// scriptedSource returns a sequence of results and records calls.
type scriptedSource struct {
	// mu protects script state and call counters across Watch goroutines.
	mu sync.Mutex

	// script is consumed one entry at a time, keeping the last entry sticky.
	script []scriptedResult

	// calls is the total number of Evaluate calls.
	calls int

	// called notifies tests that Evaluate advanced.
	called chan int

	// targets records evaluated targets in call order.
	targets []health.Target
}

// newScriptedSource returns a Source backed by script.
func newScriptedSource(script ...scriptedResult) *scriptedSource {
	return &scriptedSource{
		script: script,
		called: make(chan int, 16),
	}
}

// Evaluate returns the next scripted result and records the target.
func (s *scriptedSource) Evaluate(_ context.Context, target health.Target) (health.Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls++
	s.targets = append(s.targets, target)

	result := scriptedResult{status: health.StatusUnknown}
	if len(s.script) > 0 {
		result = s.script[0]
		if len(s.script) > 1 {
			s.script = s.script[1:]
		}
	}

	select {
	case s.called <- s.calls:
	default:
	}

	if result.err != nil {
		return health.Report{}, result.err
	}

	return testReport(target, result.status), nil
}

// callCount returns the number of recorded Evaluate calls.
func (s *scriptedSource) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls
}

// countingSource records per-target evaluations for List deduplication tests.
type countingSource struct {
	// mu protects maps because List and Watch tests may run sources concurrently.
	mu sync.Mutex

	// statuses maps each target to the status Evaluate should return.
	statuses map[health.Target]health.Status

	// errors maps targets to evaluation errors.
	errors map[health.Target]error

	// calls records evaluation counts by target.
	calls map[health.Target]int
}

// newCountingSource returns an initialized counting Source.
func newCountingSource() *countingSource {
	return &countingSource{
		statuses: make(map[health.Target]health.Status),
		errors:   make(map[health.Target]error),
		calls:    make(map[health.Target]int),
	}
}

// Evaluate records target and returns the configured status or error.
func (s *countingSource) Evaluate(_ context.Context, target health.Target) (health.Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls[target]++
	if err := s.errors[target]; err != nil {
		return health.Report{}, err
	}

	status, ok := s.statuses[target]
	if !ok {
		status = health.StatusHealthy
	}

	return testReport(target, status), nil
}

// callsFor returns how many times target has been evaluated.
func (s *countingSource) callsFor(target health.Target) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls[target]
}

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

// testReport returns a minimal valid report for target and status.
func testReport(target health.Target, status health.Status) health.Report {
	return health.Report{
		Target:   target,
		Status:   status,
		Observed: testObserved,
	}
}

// grpcCode returns the canonical gRPC code for err.
func grpcCode(err error) codes.Code {
	return status.Code(err)
}

// testClock returns a FakeClock anchored at testObserved.
func testClock() *clock.FakeClock {
	return clock.NewFakeClock(testObserved)
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

// waitForCalls waits until scriptedSource has observed at least want calls.
func waitForCalls(t *testing.T, source *scriptedSource, want int) {
	t.Helper()

	if source.callCount() >= want {
		return
	}

	deadline := time.After(testTimeout)
	for {
		select {
		case <-source.called:
			if source.callCount() >= want {
				return
			}
		case <-deadline:
			t.Fatalf("timed out waiting for source calls >= %d", want)
		}
	}
}
