// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wait

import (
	"context"
	"testing"
	"time"

	panicassert "arcoris.dev/testutil/panic"
)

func TestUntilRejectsNilOption(t *testing.T) {
	t.Parallel()

	panicassert.RequireValue(t, errNilOption, func() {
		_ = Until(context.Background(), time.Second, Satisfied, nil)
	})
}

func TestTimerStopDoesNotCloseChannel(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)
	timer.Stop()

	select {
	case _, ok := <-timer.C():
		if !ok {
			t.Fatal("timer channel was closed by Stop")
		}
	default:
	}
}

func TestTimerStopAndDrainDrainsPendingDelivery(t *testing.T) {
	t.Parallel()

	timer := NewTimer(0)
	select {
	case <-timer.C():
	case <-time.After(time.Second):
		t.Fatal("timer did not fire")
	}

	timer.StopAndDrain()
	select {
	case val := <-timer.C():
		t.Fatalf("received stale timer value %v after StopAndDrain", val)
	default:
	}
}

func TestTimerResetAfterFiredTimerDoesNotDeliverStaleValue(t *testing.T) {
	t.Parallel()

	timer := NewTimer(0)
	select {
	case <-timer.C():
	case <-time.After(time.Second):
		t.Fatal("timer did not fire")
	}

	timer.Reset(time.Hour)
	defer timer.StopAndDrain()

	select {
	case val := <-timer.C():
		t.Fatalf("received stale timer value %v after Reset", val)
	default:
	}
}

func TestTimerWaitContextStopDrainsTimer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	timer := NewTimer(0)

	err := timer.Wait(ctx)
	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)

	select {
	case val := <-timer.C():
		t.Fatalf("received timer value %v after context-stopped Wait", val)
	default:
	}
}

func TestJitterNeverShortensInterval(t *testing.T) {
	t.Parallel()

	base := time.Second
	for i := 0; i < 128; i++ {
		if got := Jitter(base, 0.5); got < base {
			t.Fatalf("Jitter(%v, 0.5) = %v, want >= base", base, got)
		}
	}
}

func TestJitterBounds(t *testing.T) {
	t.Parallel()

	base := 20 * time.Millisecond
	got := jitterWithDraw(base, 0.25, func(maxDelta time.Duration) time.Duration {
		return maxDelta
	})
	if want := 25 * time.Millisecond; got != want {
		t.Fatalf("jitter upper bound = %v, want %v", got, want)
	}
}
