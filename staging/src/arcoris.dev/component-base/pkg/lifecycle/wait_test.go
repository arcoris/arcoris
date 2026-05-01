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
	"time"
)

func TestWaitImmediateSuccess(t *testing.T) {
	t.Parallel()

	snapshot, err := NewController().Wait(nil, func(snapshot Snapshot) bool {
		return snapshot.State == StateNew
	})
	if err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if snapshot.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snapshot.State)
	}
}

func TestWaitNilPredicate(t *testing.T) {
	t.Parallel()

	_, err := NewController().Wait(context.Background(), nil)
	if err == nil {
		t.Fatal("Wait nil predicate err = nil, want error")
	}
	mustMatch(t, err, ErrInvalidWaitPredicate)
}

func TestWaitContextStops(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ctx    context.Context
		cancel context.CancelFunc
		want   error
	}{
		{
			name: "cancelled",
			want: context.Canceled,
		},
		{
			name: "deadline",
			want: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			var cancel context.CancelFunc
			if tt.want == context.Canceled {
				ctx, cancel = context.WithCancel(context.Background())
				cancel()
			} else {
				ctx, cancel = context.WithTimeout(context.Background(), 0)
				defer cancel()
			}

			_, err := NewController().Wait(ctx, func(Snapshot) bool { return false })
			if err == nil {
				t.Fatal("Wait err = nil, want context error")
			}
			mustMatch(t, err, tt.want)
		})
	}
}

func TestWaitTargetUnreachableAfterTerminal(t *testing.T) {
	t.Parallel()

	controller := NewController()
	if _, err := controller.BeginStop(); err != nil {
		t.Fatalf("BeginStop = %v", err)
	}

	snapshot, err := controller.Wait(context.Background(), func(Snapshot) bool { return false })
	if err == nil {
		t.Fatal("Wait err = nil, want unreachable")
	}
	mustMatch(t, err, ErrWaitTargetUnreachable)
	if snapshot.State != StateStopped {
		t.Fatalf("snapshot.State = %s, want stopped", snapshot.State)
	}
}

func TestWaitState(t *testing.T) {
	t.Parallel()

	controller := NewController()
	snapshot, err := controller.WaitState(context.Background(), StateNew)
	if err != nil {
		t.Fatalf("WaitState current = %v", err)
	}
	if snapshot.State != StateNew {
		t.Fatalf("snapshot.State = %s, want new", snapshot.State)
	}

	_, err = controller.WaitState(context.Background(), State(99))
	if err == nil {
		t.Fatal("WaitState invalid target err = nil, want error")
	}
	mustMatch(t, err, ErrInvalidWaitTarget)

	if _, err := controller.BeginStart(); err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	_, err = controller.WaitState(context.Background(), StateNew)
	if err == nil {
		t.Fatal("WaitState backward err = nil, want unreachable")
	}
	mustMatch(t, err, ErrWaitTargetUnreachable)
}

func TestWaitStateWaitsUntilTarget(t *testing.T) {
	t.Parallel()

	controller := NewController()
	results := make(chan Snapshot, 1)
	errs := make(chan error, 1)

	go func() {
		snapshot, err := controller.WaitState(context.Background(), StateRunning)
		if err != nil {
			errs <- err
			return
		}
		results <- snapshot
	}()

	if _, err := controller.BeginStart(); err != nil {
		t.Fatalf("BeginStart = %v", err)
	}
	if _, err := controller.MarkRunning(); err != nil {
		t.Fatalf("MarkRunning = %v", err)
	}

	select {
	case err := <-errs:
		t.Fatalf("WaitState err = %v, want nil", err)
	default:
	}
	if got := mustReceiveSnapshot(t, results); got.State != StateRunning {
		t.Fatalf("snapshot.State = %s, want running", got.State)
	}
}

func TestWaitTerminal(t *testing.T) {
	t.Parallel()

	t.Run("stopped", func(t *testing.T) {
		t.Parallel()

		controller := NewController()
		if _, err := controller.BeginStop(); err != nil {
			t.Fatalf("BeginStop = %v", err)
		}
		snapshot, err := controller.WaitTerminal(context.Background())
		if err != nil {
			t.Fatalf("WaitTerminal = %v, want nil", err)
		}
		if snapshot.State != StateStopped {
			t.Fatalf("snapshot.State = %s, want stopped", snapshot.State)
		}
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("failed")
		controller := NewController()
		if _, err := controller.BeginStart(); err != nil {
			t.Fatalf("BeginStart = %v", err)
		}
		if _, err := controller.MarkFailed(cause); err != nil {
			t.Fatalf("MarkFailed = %v", err)
		}
		snapshot, err := controller.WaitTerminal(context.Background())
		if err != nil {
			t.Fatalf("WaitTerminal = %v, want nil", err)
		}
		if snapshot.State != StateFailed || snapshot.FailureCause != cause {
			t.Fatalf("snapshot = %+v, want failed with cause", snapshot)
		}
	})
}

func TestDoneStableAndClosesOnce(t *testing.T) {
	t.Parallel()

	controller := NewController()
	first := controller.Done()
	second := controller.Done()
	if first != second {
		t.Fatal("Done returned different channels")
	}
	mustNotSignalClosed(t, first)

	if _, err := controller.BeginStop(); err != nil {
		t.Fatalf("BeginStop = %v", err)
	}
	mustSignalClosed(t, first)
	mustSignalClosed(t, second)
}

func TestWaitReturnsLatestSnapshotOnContextCancel(t *testing.T) {
	t.Parallel()

	controller := NewController()
	if _, err := controller.BeginStart(); err != nil {
		t.Fatalf("BeginStart = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	snapshot, err := controller.Wait(ctx, func(Snapshot) bool { return false })
	if err == nil {
		t.Fatal("Wait err = nil, want canceled")
	}
	mustMatch(t, err, context.Canceled)
	if snapshot.State != StateStarting {
		t.Fatalf("snapshot.State = %s, want starting", snapshot.State)
	}
}

func TestWaitDeadlineExceededDuringWait(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	_, err := NewController().Wait(ctx, func(Snapshot) bool { return false })
	if err == nil {
		t.Fatal("Wait err = nil, want deadline")
	}
	mustMatch(t, err, context.DeadlineExceeded)
}
