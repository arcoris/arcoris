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
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestWatchRejectsNilRequest(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))
	err := server.Watch(nil, newWatchStream())
	if grpcCode(err) != codes.InvalidArgument {
		t.Fatalf("Watch(nil) code = %s, want InvalidArgument", grpcCode(err))
	}
}

func TestWatchRejectsNilServerOrStream(t *testing.T) {
	t.Parallel()

	var server *Server
	if err := server.Watch(&healthpb.HealthCheckRequest{}, newWatchStream()); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch(nil server) code = %s, want Canceled", grpcCode(err))
	}

	server = mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))
	if err := server.Watch(&healthpb.HealthCheckRequest{}, nil); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch(nil stream) code = %s, want Canceled", grpcCode(err))
	}
}

func TestWatchUnknownServiceSendsServiceUnknown(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{Service: "missing"}, stream)
	}()

	if got := mustReceiveStatus(t, stream); got != healthpb.HealthCheckResponse_SERVICE_UNKNOWN {
		t.Fatalf("initial status = %s, want SERVICE_UNKNOWN", got)
	}

	stream.cancel()
	if err := waitForWatchDone(t, done); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
}

func TestWatchKnownServiceSendsInitialStatus(t *testing.T) {
	t.Parallel()

	clk := testClock()
	server := mustNewServer(
		t,
		healthtest.NewSequenceSource(health.TargetReady, healthtest.HealthyReport(health.TargetReady)),
		WithClock(clk),
		WithWatchInterval(time.Second),
	)
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{}, stream)
	}()

	if got := mustReceiveStatus(t, stream); got != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("initial status = %s, want SERVING", got)
	}

	stream.cancel()
	if err := waitForWatchDone(t, done); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
}

func TestWatchSendsChangedStatus(t *testing.T) {
	t.Parallel()

	clk := testClock()
	source := healthtest.NewSequenceSource(
		health.TargetReady,
		healthtest.HealthyReport(health.TargetReady),
		healthtest.UnhealthyReport(health.TargetReady),
	)
	server := mustNewServer(t, source, WithClock(clk), WithWatchInterval(time.Second))
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{}, stream)
	}()

	if got := mustReceiveStatus(t, stream); got != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("initial status = %s, want SERVING", got)
	}
	if got := stepClockUntilStatus(t, clk, stream, time.Second); got != healthpb.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("changed status = %s, want NOT_SERVING", got)
	}

	stream.cancel()
	if err := waitForWatchDone(t, done); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
}

func TestWatchDoesNotSendDuplicateStatus(t *testing.T) {
	t.Parallel()

	clk := testClock()
	recording := &recordingClock{Clock: clk, created: make(chan time.Duration, 1)}
	source := healthtest.NewSequenceSource(
		health.TargetReady,
		healthtest.HealthyReport(health.TargetReady),
		healthtest.HealthyReport(health.TargetReady),
		healthtest.UnhealthyReport(health.TargetReady),
	)
	server := mustNewServer(t, source, WithClock(recording), WithWatchInterval(time.Second))
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{}, stream)
	}()

	if got := mustReceiveStatus(t, stream); got != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("initial status = %s, want SERVING", got)
	}
	waitForTickerInterval(t, recording, time.Second)

	clk.Step(time.Second)
	waitForSourceCalls(t, source, health.TargetReady, 2)
	assertNoWatchStatus(t, stream)

	clk.Step(time.Second)
	if got := mustReceiveStatus(t, stream); got != healthpb.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("changed status = %s, want NOT_SERVING", got)
	}

	stream.cancel()
	if err := waitForWatchDone(t, done); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
}

func TestWatchSourceErrorMapsToUnknown(t *testing.T) {
	t.Parallel()

	raw := errors.New("raw password=secret source error")
	server := mustNewServer(
		t,
		healthtest.NewErrorSource(raw),
		WithClock(testClock()),
		WithWatchInterval(time.Second),
	)
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{}, stream)
	}()

	if got := mustReceiveStatus(t, stream); got != healthpb.HealthCheckResponse_UNKNOWN {
		t.Fatalf("initial status = %s, want UNKNOWN", got)
	}

	stream.cancel()
	err := waitForWatchDone(t, done)
	if grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
	if strings.Contains(err.Error(), "password=secret") {
		t.Fatalf("Watch() leaked raw source error: %v", err)
	}
}

func TestWatchStreamCancellationStopsWatch(t *testing.T) {
	t.Parallel()

	server := mustNewServer(
		t,
		healthtest.NewSequenceSource(health.TargetReady, healthtest.HealthyReport(health.TargetReady)),
		WithClock(testClock()),
		WithWatchInterval(time.Second),
	)
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{}, stream)
	}()

	_ = mustReceiveStatus(t, stream)
	stream.cancel()

	err := waitForWatchDone(t, done)
	if grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
}

func TestWatchSendErrorReturnsCanceled(t *testing.T) {
	t.Parallel()

	stream := newWatchStream()
	stream.sendErr = errors.New("raw stream transport failed")
	server := mustNewServer(t, healthtest.NewStaticSource(healthtest.HealthyReport(health.TargetReady)))

	err := server.Watch(&healthpb.HealthCheckRequest{}, stream)
	if grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
	if !strings.Contains(err.Error(), watchEndedMessage) {
		t.Fatalf("Watch() = %v, want generic watch-ended message", err)
	}
	if strings.Contains(err.Error(), "raw stream transport failed") {
		t.Fatalf("Watch() leaked raw send error: %v", err)
	}
}

func TestWatchUsesConfiguredInterval(t *testing.T) {
	t.Parallel()

	clk := &recordingClock{Clock: testClock(), created: make(chan time.Duration, 1)}
	server := mustNewServer(
		t,
		healthtest.NewSequenceSource(health.TargetReady, healthtest.HealthyReport(health.TargetReady)),
		WithClock(clk),
		WithWatchInterval(2*time.Second),
	)
	stream := newWatchStream()
	done := make(chan error, 1)

	go func() {
		done <- server.Watch(&healthpb.HealthCheckRequest{}, stream)
	}()

	_ = mustReceiveStatus(t, stream)
	select {
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for ticker creation")
	case interval := <-clk.created:
		if interval != 2*time.Second {
			t.Fatalf("ticker interval = %s, want 2s", interval)
		}
	}

	stream.cancel()
	if err := waitForWatchDone(t, done); grpcCode(err) != codes.Canceled {
		t.Fatalf("Watch() code = %s, want Canceled", grpcCode(err))
	}
}

// recordingClock wraps a test clock and records ticker intervals.
type recordingClock struct {
	// Clock is the underlying clock used to create real fake tickers.
	clock.Clock

	// mu protects non-blocking writes to created.
	mu sync.Mutex

	// created receives each interval passed to NewTicker.
	created chan time.Duration
}

// NewTicker records d and delegates ticker creation to the wrapped clock.
func (c *recordingClock) NewTicker(d time.Duration) clock.Ticker {
	c.mu.Lock()
	select {
	case c.created <- d:
	default:
	}
	c.mu.Unlock()

	return c.Clock.NewTicker(d)
}

// stepClockUntilStatus advances clk until Watch sends a status.
func stepClockUntilStatus(
	t *testing.T,
	clk *clock.FakeClock,
	stream *watchStream,
	step time.Duration,
) healthpb.HealthCheckResponse_ServingStatus {
	t.Helper()

	deadline := time.After(testTimeout)
	for {
		clk.Step(step)
		select {
		case resp := <-stream.sent:
			return resp.GetStatus()
		case <-deadline:
			t.Fatal("timed out waiting for watch status after clock step")
			return healthpb.HealthCheckResponse_UNKNOWN
		default:
		}
	}
}

// waitForTickerInterval waits for recordingClock to observe want.
func waitForTickerInterval(t *testing.T, clk *recordingClock, want time.Duration) {
	t.Helper()

	select {
	case interval := <-clk.created:
		if interval != want {
			t.Fatalf("ticker interval = %s, want %s", interval, want)
		}
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for ticker creation")
	}
}
