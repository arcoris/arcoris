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
)

func TestGroupJoinErrorsAreOrderedBySubmission(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background(), WithCancelOnError(false))
	firstErr := errors.New("first")
	secondErr := errors.New("second")

	releaseFirst := make(chan struct{})
	secondReturned := make(chan struct{})

	group.Go("first", func(ctx context.Context) error {
		<-releaseFirst
		return firstErr
	})
	group.Go("second", func(ctx context.Context) error {
		close(secondReturned)
		return secondErr
	})

	mustClose(t, secondReturned)
	close(releaseFirst)

	err := group.Wait()
	if _, ok := err.(interface{ Unwrap() []error }); !ok {
		t.Fatalf("Wait error type %T does not expose joined errors", err)
	}

	taskErrs := TaskErrors(err)
	if len(taskErrs) != 2 {
		t.Fatalf("TaskErrors len = %d, want 2", len(taskErrs))
	}
	if taskErrs[0].Name != "first" || taskErrs[1].Name != "second" {
		t.Fatalf("TaskErrors order = %+v, want first then second", taskErrs)
	}
}

func TestGroupJoinSingleErrorStillReturnsJoinedTaskError(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	want := errors.New("boom")

	group.Go("worker", func(ctx context.Context) error {
		return want
	})

	err := group.Wait()
	if _, ok := err.(interface{ Unwrap() []error }); !ok {
		t.Fatalf("Wait error type %T does not expose joined errors", err)
	}
	if !errors.Is(err, want) {
		t.Fatalf("Wait error = %v, want original error match", err)
	}

	var taskErr TaskError
	if !errors.As(err, &taskErr) {
		t.Fatal("Wait error does not expose TaskError")
	}

	taskErrs := TaskErrors(err)
	if len(taskErrs) != 1 {
		t.Fatalf("TaskErrors len = %d, want 1", len(taskErrs))
	}
	if taskErrs[0].Name != "worker" {
		t.Fatalf("TaskError name = %q, want worker", taskErrs[0].Name)
	}
}

func TestGroupErrorModeFirstReturnsFirstObservedError(t *testing.T) {
	t.Parallel()

	group := NewGroup(
		context.Background(),
		WithCancelOnError(false),
		WithErrorMode(ErrorModeFirst),
	)
	firstErr := errors.New("first")
	secondErr := errors.New("second")

	releaseFirst := make(chan struct{})
	secondReturned := make(chan struct{})

	group.Go("first", func(ctx context.Context) error {
		<-releaseFirst
		return firstErr
	})
	group.Go("second", func(ctx context.Context) error {
		close(secondReturned)
		return secondErr
	})

	mustClose(t, secondReturned)
	close(releaseFirst)

	err := group.Wait()
	taskErrs := TaskErrors(err)
	if len(taskErrs) != 1 {
		t.Fatalf("TaskErrors len = %d, want 1", len(taskErrs))
	}
	if taskErrs[0].Name != "second" {
		t.Fatalf("first observed task = %q, want second", taskErrs[0].Name)
	}
}

func TestGroupConcurrentTaskErrorsAreRaceSafe(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background(), WithCancelOnError(false))
	start := make(chan struct{})

	for _, name := range []string{"a", "b", "c", "d"} {
		name := name
		group.Go(name, func(ctx context.Context) error {
			<-start
			return errors.New(name)
		})
	}

	close(start)

	err := group.Wait()
	taskErrs := TaskErrors(err)
	if len(taskErrs) != 4 {
		t.Fatalf("TaskErrors len = %d, want 4", len(taskErrs))
	}
}

func TestGroupTaskErrorCancellationUsesTaskErrorCause(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	want := errors.New("boom")

	group.Go("worker", func(ctx context.Context) error {
		return want
	})

	mustClose(t, group.Done())

	var taskErr TaskError
	if !errors.As(context.Cause(group.Context()), &taskErr) {
		t.Fatalf("context cause = %v, want TaskError", context.Cause(group.Context()))
	}
	if taskErr.Name != "worker" {
		t.Fatalf("TaskError name = %q, want worker", taskErr.Name)
	}
	if !errors.Is(taskErr, want) {
		t.Fatalf("TaskError = %v, want original error match", taskErr)
	}

	_ = group.Wait()
}
