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

package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	panicassert "arcoris.dev/testutil/panic"
)

func TestGroupConcurrentGoAndWaitDoesNotLoseReservedTasks(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	start := make(chan struct{})
	releaseTasks := make(chan struct{})
	var attempted sync.WaitGroup
	var started atomic.Int64
	var completed atomic.Int64

	const tasks = 64
	attempted.Add(tasks)
	for i := 0; i < tasks; i++ {
		i := i
		go func() {
			defer attempted.Done()
			<-start
			defer func() {
				if recovered := recover(); recovered != nil && recovered != errGroupClosed {
					panic(recovered)
				}
			}()
			group.Go(fmt.Sprintf("task-%02d", i), func(context.Context) error {
				started.Add(1)
				<-releaseTasks
				completed.Add(1)
				return nil
			})
		}()
	}

	waitDone := make(chan error, 1)
	go func() {
		<-start
		waitDone <- group.Wait()
	}()

	close(start)
	attempted.Wait()
	close(releaseTasks)

	if err := <-waitDone; err != nil {
		t.Fatalf("Wait error = %v, want nil", err)
	}
	if got, want := completed.Load(), started.Load(); got != want {
		t.Fatalf("completed tasks = %d, want started task count %d", got, want)
	}
}

func TestGroupConcurrentCancelAndWaitIsRaceFree(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background())
	release := make(chan struct{})
	group.Go("worker", func(ctx context.Context) error {
		<-release
		return nil
	})

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		group.Cancel(nil)
	}()
	go func() {
		defer wg.Done()
		<-start
		_ = group.Wait()
	}()

	close(start)
	close(release)
	wg.Wait()

	if err := group.Wait(); err != nil {
		t.Fatalf("cached Wait error = %v, want nil", err)
	}
}

func TestGroupConcurrentTaskErrorsAreAllRecorded(t *testing.T) {
	t.Parallel()

	group := NewGroup(context.Background(), WithCancelOnError(false))
	start := make(chan struct{})
	const tasks = 16

	for i := 0; i < tasks; i++ {
		i := i
		group.Go(fmt.Sprintf("task-%02d", i), func(context.Context) error {
			<-start
			return fmt.Errorf("failure-%02d", i)
		})
	}

	close(start)
	err := group.Wait()
	if got := len(TaskErrors(err)); got != tasks {
		t.Fatalf("TaskErrors len = %d, want %d", got, tasks)
	}
}

func TestTaskErrorsWalksWrappedTaskErrors(t *testing.T) {
	t.Parallel()

	want := TaskError{Name: "worker", Err: errors.New("boom")}
	err := fmt.Errorf("outer: %w", want)

	got := TaskErrors(err)
	if len(got) != 1 {
		t.Fatalf("TaskErrors len = %d, want 1", len(got))
	}
	if got[0].Name != want.Name || !errors.Is(got[0], want.Err) {
		t.Fatalf("TaskErrors()[0] = %+v, want %+v", got[0], want)
	}
}

func TestGroupTaskPanicIsNotRecovered(t *testing.T) {
	if os.Getenv("ARCORIS_RUNTIME_RUN_PANIC_SUBPROCESS") == "1" {
		group := NewGroup(context.Background())
		group.Go("panic", func(context.Context) error {
			panic("task panic")
		})
		_ = group.Wait()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestGroupTaskPanicIsNotRecovered$")
	cmd.Env = append(os.Environ(), "ARCORIS_RUNTIME_RUN_PANIC_SUBPROCESS=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("subprocess succeeded, want task panic to crash process")
	}
	if !strings.Contains(string(out), "task panic") {
		t.Fatalf("subprocess output = %s, want task panic", out)
	}
}

func TestGroupZeroValuePanics(t *testing.T) {
	t.Parallel()

	var group Group
	tests := []struct {
		name string
		call func()
	}{
		{name: "Go", call: func() { group.Go("task", func(context.Context) error { return nil }) }},
		{name: "Context", call: func() { _ = group.Context() }},
		{name: "Done", call: func() { _ = group.Done() }},
		{name: "Cancel", call: func() { group.Cancel(nil) }},
		{name: "Wait", call: func() { _ = group.Wait() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			panicassert.RequireMessage(t, errUninitializedGroup, tt.call)
		})
	}
}
