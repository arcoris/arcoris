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

package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/clock"
)

func TestWaitDelayReturnsImmediatelyForZeroDelay(t *testing.T) {
	err := waitDelay(context.Background(), clock.RealClock{}, 0)
	if err != nil {
		t.Fatalf("waitDelay returned error: %v", err)
	}
}

func TestWaitDelayPanicsOnNegativeDelay(t *testing.T) {
	expectPanic(t, panicNegativeBackoffDelay, func() {
		_ = waitDelay(context.Background(), clock.RealClock{}, -time.Nanosecond)
	})
}

func TestWaitDelayReturnsInterruptedWhenContextAlreadyStopped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitDelay(ctx, clock.RealClock{}, time.Second)
	if err == nil {
		t.Fatalf("waitDelay returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("waitDelay error does not match ErrInterrupted: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitDelay error does not preserve context.Canceled: %v", err)
	}
}

func TestWaitDelayReturnsInterruptedWhenContextStopsDuringDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- waitDelay(ctx, clock.RealClock{}, time.Hour)
	}()

	cancel()

	err := <-done
	if err == nil {
		t.Fatalf("waitDelay returned nil")
	}
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("waitDelay error does not match ErrInterrupted: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitDelay error does not preserve context.Canceled: %v", err)
	}
}
