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
	"sync"
)

// Group owns a single context-first runtime task scope.
//
// A Group is created with NewGroup, starts named tasks with Go, exposes the
// shared task context through Context and Done, can be cancelled by its owner
// with Cancel, and is joined through Wait. Group values must not be copied after
// construction.
type Group struct {
	noCopy noCopy

	ctx    context.Context
	cancel context.CancelCauseFunc

	wg sync.WaitGroup

	mu      sync.Mutex
	closed  bool
	nextSeq uint64
	names   map[string]struct{}
	errs    []taskErrorRecord

	// taskCancelSet records that fail-fast cancellation has already been claimed
	// by the first task error recorded under mu. Later task errors are still
	// appended for Wait, but they must not race to overwrite the context cause.
	taskCancelSet bool

	config groupConfig

	waitOnce sync.Once
	waitErr  error
}

// NewGroup creates a Group derived from parent.
//
// NewGroup panics when parent is nil or when options produce invalid
// configuration. The returned Group is single-use and must not be copied after
// construction.
func NewGroup(parent context.Context, opts ...GroupOption) *Group {
	requireContext(parent, errNilGroupParent)

	ctx, cancel := context.WithCancelCause(parent)
	return &Group{
		ctx:    ctx,
		cancel: cancel,
		names:  make(map[string]struct{}),
		config: newGroupConfig(opts...),
	}
}

// Go starts task in a new goroutine under name.
//
// The task receives the Group context. If task returns a non-nil error, the
// error is recorded as a TaskError. When cancel-on-error is enabled, the first
// non-nil task error also cancels the Group context.
//
// Go panics when the Group is nil or uninitialized, name is invalid, task is nil,
// name has already been used, or the Group has been closed by Wait, Cancel, or
// fail-fast cancellation.
func (g *Group) Go(name string, task Task) {
	requireGroup(g)
	requireTaskName(name)
	requireTask(task)

	seq := g.reserveTask(name)

	go func() {
		defer g.wg.Done()

		if err := task(g.ctx); err != nil {
			g.recordTaskError(seq, name, err)
		}
	}()
}

// Context returns the context shared by all tasks in the Group.
//
// The returned context is stable for the lifetime of the Group. It is cancelled
// by owner Cancel, by the first task error when cancel-on-error is enabled, by
// the parent context, or by Wait after all started tasks have completed.
func (g *Group) Context() context.Context {
	requireGroup(g)
	return g.ctx
}

// Done returns the Group context Done channel.
//
// Done is equivalent to Context().Done and is provided for owners that only need
// to select on group cancellation.
func (g *Group) Done() <-chan struct{} {
	requireGroup(g)
	return g.ctx.Done()
}

// Cancel cancels the Group context and closes the Group for new task
// submissions.
//
// Cancel is idempotent. A nil cause is normalized to context.Canceled. Cancel
// does not record a task error and does not wait for running tasks.
func (g *Group) Cancel(cause error) {
	requireGroup(g)

	if cause == nil {
		cause = context.Canceled
	}

	g.close()
	g.cancel(cause)
}

// Wait closes the Group for new task submissions and waits for all started tasks
// to finish.
//
// Wait is idempotent. The first call waits for all tasks, releases the group
// context, builds the configured task error result, and caches that result. Later
// calls return the cached result.
func (g *Group) Wait() error {
	requireGroup(g)

	g.waitOnce.Do(func() {
		g.close()
		g.wg.Wait()
		g.cancel(context.Canceled)
		g.waitErr = g.buildWaitError()
	})

	return g.waitErr
}
