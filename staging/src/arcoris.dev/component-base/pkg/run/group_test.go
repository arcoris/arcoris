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
	"testing"
	"time"
)

func TestNewGroupRejectsNilParent(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilGroupParent, func() {
		NewGroup(nil)
	})
}

func TestGroupContextAndDoneAreStable(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	first := group.Context()

	if first == nil {
		t.Fatal("Context returned nil")
	}
	if group.Context() != first {
		t.Fatal("Context did not return a stable context")
	}
	if group.Done() != first.Done() {
		t.Fatal("Done did not return Context().Done()")
	}
}

func TestNewGroupCreatesChildContext(t *testing.T) {
	t.Parallel()

	parent, cancel := context.WithCancelCause(context.Background())
	group := NewGroup(parent)

	if group.Context() == parent {
		t.Fatal("NewGroup returned parent context directly")
	}

	want := errors.New("parent stop")
	cancel(want)
	mustClose(t, group.Done())

	if !errors.Is(context.Cause(group.Context()), want) {
		t.Fatalf("context cause = %v, want %v", context.Cause(group.Context()), want)
	}
}

func TestGroupGoStartsTaskAndWaitReturnsNil(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	called := make(chan struct{})

	group.Go("worker", func(ctx context.Context) error {
		close(called)
		return nil
	})

	if err := group.Wait(); err != nil {
		t.Fatalf("Wait error = %v", err)
	}
	mustClose(t, called)
}

func TestGroupTaskReceivesGroupContext(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())

	group.Go("worker", func(ctx context.Context) error {
		if ctx != group.Context() {
			t.Fatal("task received a different context")
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		t.Fatalf("Wait error = %v", err)
	}
}

func TestGroupTaskErrorCancelsContextByDefault(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	want := errors.New("boom")

	group.Go("worker", func(ctx context.Context) error {
		return want
	})

	err := group.Wait()
	if !errors.Is(err, want) {
		t.Fatalf("Wait error = %v, want %v", err, want)
	}
	if !errors.Is(context.Cause(group.Context()), want) {
		t.Fatalf("context cause = %v, want %v", context.Cause(group.Context()), want)
	}

	var taskErr TaskError
	if !errors.As(err, &taskErr) {
		t.Fatal("Wait error does not contain TaskError")
	}
	if taskErr.Name != "worker" {
		t.Fatalf("TaskError name = %q, want worker", taskErr.Name)
	}
}

func TestGroupWithCancelOnErrorDisabledDoesNotCancelOnTaskError(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background(), WithCancelOnError(false))
	release := make(chan struct{})
	failed := make(chan struct{})
	cancelled := make(chan struct{})
	want := errors.New("boom")

	group.Go("failing", func(ctx context.Context) error {
		close(failed)
		return want
	})
	group.Go("observer", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			close(cancelled)
			return nil
		case <-release:
			return nil
		}
	})

	mustClose(t, failed)
	mustNotCloseNow(t, cancelled)

	close(release)

	if err := group.Wait(); !errors.Is(err, want) {
		t.Fatalf("Wait error = %v, want %v", err, want)
	}
}

func TestGroupCancelCancelsContextWithoutTaskError(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())

	group.Go("worker", func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})

	group.Cancel(nil)

	if err := group.Wait(); err != nil {
		t.Fatalf("Wait error = %v, want nil", err)
	}
	if taskErrs := TaskErrors(group.Wait()); len(taskErrs) != 0 {
		t.Fatalf("Cancel recorded task errors: %+v", taskErrs)
	}
	if !errors.Is(context.Cause(group.Context()), context.Canceled) {
		t.Fatalf("context cause = %v, want context.Canceled", context.Cause(group.Context()))
	}
}

func TestGroupCancelWithCause(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	want := errors.New("owner stop")

	group.Cancel(want)
	if !errors.Is(context.Cause(group.Context()), want) {
		t.Fatalf("context cause = %v, want %v", context.Cause(group.Context()), want)
	}
}

func TestGroupWaitIsIdempotent(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	want := errors.New("boom")

	group.Go("worker", func(ctx context.Context) error {
		return want
	})

	first := group.Wait()
	second := group.Wait()

	if first == nil || second == nil {
		t.Fatal("Wait returned nil, want error")
	}
	if first != second {
		t.Fatalf("Wait did not return cached error: first=%p second=%p", first, second)
	}
}

func TestGroupWaitWaitsForAllTasksAfterFirstError(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	fail := make(chan struct{})
	release := make(chan struct{})
	secondDone := make(chan struct{})
	want := errors.New("boom")

	group.Go("failing", func(ctx context.Context) error {
		<-fail
		return want
	})
	group.Go("cleanup", func(ctx context.Context) error {
		<-release
		close(secondDone)
		return nil
	})

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- group.Wait()
	}()

	close(fail)
	mustClose(t, group.Done())
	mustNotCloseNow(t, waitDone)

	close(release)
	mustClose(t, secondDone)

	select {
	case err := <-waitDone:
		if !errors.Is(err, want) {
			t.Fatalf("Wait error = %v, want %v", err, want)
		}
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for Wait")
	}
}

func TestGroupRejectsInvalidGoInputs(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())

	mustPanicWith(t, errEmptyTaskName, func() {
		group.Go("", func(ctx context.Context) error { return nil })
	})
	mustPanicWith(t, errUntrimmedTaskName, func() {
		group.Go(" worker", func(ctx context.Context) error { return nil })
	})
	mustPanicWith(t, errNilTask, func() {
		group.Go("worker", nil)
	})
}

func TestGroupRejectsDuplicateTaskName(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	release := make(chan struct{})

	group.Go("worker", func(ctx context.Context) error {
		<-release
		return nil
	})

	mustPanicWith(t, errDuplicateTaskName, func() {
		group.Go("worker", func(ctx context.Context) error { return nil })
	})

	close(release)
	if err := group.Wait(); err != nil {
		t.Fatalf("Wait error = %v", err)
	}
}

func TestGroupRejectsGoAfterWait(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	if err := group.Wait(); err != nil {
		t.Fatalf("Wait error = %v", err)
	}

	mustPanicWith(t, errGroupClosed, func() {
		group.Go("worker", func(ctx context.Context) error { return nil })
	})
}

func TestGroupRejectsGoAfterCancel(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	group.Cancel(nil)

	mustPanicWith(t, errGroupClosed, func() {
		group.Go("worker", func(ctx context.Context) error { return nil })
	})
}

func TestGroupRejectsNilAndUninitializedReceiver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(*Group)
	}{
		{name: "Go", fn: func(g *Group) { g.Go("task", func(ctx context.Context) error { return nil }) }},
		{name: "Context", fn: func(g *Group) { g.Context() }},
		{name: "Done", fn: func(g *Group) { g.Done() }},
		{name: "Cancel", fn: func(g *Group) { g.Cancel(nil) }},
		{name: "Wait", fn: func(g *Group) { g.Wait() }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("nil "+tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilGroup, func() {
				tt.fn(nil)
			})
		})
		t.Run("zero "+tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errUninitializedGroup, func() {
				var group Group
				tt.fn(&group)
			})
		})
	}
}
