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

package lifecycle

import (
	"context"
	"errors"
	"testing"
)

func TestWaitTerminalReturnsStoppedSnapshot(t *testing.T) {
	t.Parallel()

	controller := NewController()
	_, _ = controller.BeginStop()
	snap, err := controller.WaitTerminal(context.Background())
	if err != nil {
		t.Fatalf("WaitTerminal = %v, want nil", err)
	}
	if snap.State != StateStopped {
		t.Fatalf("snapshot.State = %s, want stopped", snap.State)
	}
}

func TestWaitTerminalReturnsFailedSnapshotWithCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	controller := NewController()
	_, _ = controller.BeginStart()
	_, _ = controller.MarkFailed(cause)
	snap, err := controller.WaitTerminal(context.Background())
	if err != nil {
		t.Fatalf("WaitTerminal = %v, want nil", err)
	}
	if snap.State != StateFailed || snap.FailureCause != cause {
		t.Fatalf("snapshot = %+v, want failed with cause", snap)
	}
}

func TestWaitTerminalWaitsUntilTerminalCommitted(t *testing.T) {
	t.Parallel()

	// Terminal means either stopped or failed; it does not imply the stop path
	// succeeded, only that the lifecycle instance has ended.
	controller := NewController()
	results := make(chan Snapshot, 1)
	errs := make(chan error, 1)
	go func() {
		snap, err := controller.WaitTerminal(context.Background())
		if err != nil {
			errs <- err
			return
		}
		results <- snap
	}()

	_, _ = controller.BeginStart()
	_, _ = controller.BeginStop()
	_, _ = controller.MarkStopped()
	select {
	case err := <-errs:
		t.Fatalf("WaitTerminal err = %v, want nil", err)
	default:
	}
	if got := mustReceiveSnapshot(t, results); got.State != StateStopped {
		t.Fatalf("snapshot.State = %s, want stopped", got.State)
	}
}

func TestWaitTerminalContextCancellationAndDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewController().WaitTerminal(ctx)
	if err == nil {
		t.Fatal("WaitTerminal cancel err = nil, want canceled")
	}
	mustMatch(t, err, context.Canceled)

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), testTime)
	defer deadlineCancel()
	_, deadlineErr := NewController().WaitTerminal(deadlineCtx)
	if deadlineErr == nil {
		t.Fatal("WaitTerminal deadline err = nil, want deadline exceeded")
	}
	mustMatch(t, deadlineErr, context.DeadlineExceeded)
}
