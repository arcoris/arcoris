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

package run

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGroupClosesForSubmissionsAfterWaitBegins(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	started := make(chan struct{})
	release := make(chan struct{})

	group.Go("blocker", func(ctx context.Context) error {
		close(started)
		<-release
		return nil
	})
	mustClose(t, started)

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- group.Wait()
	}()

	// Wait closes submission under the same mutex used by reserveTask before it
	// waits for tasks. This is the invariant that prevents WaitGroup.Add from
	// racing with Wait.
	waitGroupClosed(t, group)
	mustPanicWith(t, errGroupClosed, func() {
		group.Go("late", func(ctx context.Context) error { return nil })
	})

	close(release)
	select {
	case err := <-waitDone:
		if err != nil {
			t.Fatalf("Wait error = %v, want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for Wait")
	}
}

func TestGroupClosesForSubmissionsAfterCancel(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	group.Cancel(nil)

	if !groupClosed(group) {
		t.Fatal("Cancel did not close group")
	}
	mustPanicWith(t, errGroupClosed, func() {
		group.Go("late", func(ctx context.Context) error { return nil })
	})
}

func TestGroupClosesForSubmissionsAfterFailFastTaskError(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	want := errors.New("boom")

	group.Go("failing", func(ctx context.Context) error {
		return want
	})

	mustClose(t, group.Done())
	mustPanicWith(t, errGroupClosed, func() {
		group.Go("late", func(ctx context.Context) error { return nil })
	})

	if err := group.Wait(); !errors.Is(err, want) {
		t.Fatalf("Wait error = %v, want %v", err, want)
	}
}

func TestGroupTaskReservationIsDeterministic(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())

	// This test exercises the reservation primitive directly because it owns the
	// submission sequence invariant. reserveTask increments the internal
	// WaitGroup before a goroutine is started, so the test manually balances that
	// accounting after inspecting the reserved names and sequence numbers. Public
	// Group.Go and Wait behavior is covered by the joined-error ordering tests.
	first := group.reserveTask("first")
	second := group.reserveTask("second")
	group.wg.Done()
	group.wg.Done()

	if first != 0 || second != 1 {
		t.Fatalf("reservation sequence = (%d, %d), want (0, 1)", first, second)
	}
	if _, ok := group.names["first"]; !ok {
		t.Fatal("first task name was not reserved")
	}
	if _, ok := group.names["second"]; !ok {
		t.Fatal("second task name was not reserved")
	}
}

func TestGroupGoMayBeCalledConcurrentlyBeforeClose(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	start := make(chan struct{})
	var submits sync.WaitGroup

	for i := 0; i < 8; i++ {
		submits.Add(1)
		go func() {
			defer submits.Done()
			<-start
			group.Go(fmt.Sprintf("task-%d", i), func(ctx context.Context) error {
				return nil
			})
		}()
	}

	close(start)
	submits.Wait()

	if err := group.Wait(); err != nil {
		t.Fatalf("Wait error = %v, want nil", err)
	}
}

func groupClosed(group *Group) bool {
	group.mu.Lock()
	defer group.mu.Unlock()

	return group.closed
}
