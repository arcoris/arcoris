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

package objectcontroller

import (
	"context"
	"sync"
)

// Run starts controller workers and blocks until they all exit.
//
// Run panics when ctx is nil. A nil Controller returns ErrInvalidController.
// Concurrent Run calls on one Controller return ErrAlreadyRunning. Sequential
// Run calls are allowed after the previous Run has returned.
//
// Run returns nil when all workers exit because the queue is shut down and
// drained. It returns ctx.Err when context cancellation stops workers, or the
// first fatal worker error. After the first fatal error, Run cancels the worker
// context and waits for every worker before returning. Run observes queue
// shutdown but never initiates it.
//
// Run does not recover reconciliation panics. A panic from a worker goroutine
// remains a panic in that goroutine rather than a returned error.
func (c *Controller) Run(ctx context.Context) error {
	if c == nil {
		return ErrInvalidController
	}
	if ctx == nil {
		panic("objectcontroller: nil context")
	}
	if err := c.startRun(); err != nil {
		return err
	}
	defer c.finishRun()

	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, c.workers)
	var wg sync.WaitGroup
	wg.Add(c.workers)
	for range c.workers {
		go func() {
			defer wg.Done()
			if err := c.runWorker(workerCtx); err != nil {
				errs <- err
				cancel()
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errs:
		cancel()
		<-done
		return err
	case <-done:
		select {
		case err := <-errs:
			return err
		default:
			return nil
		}
	}
}
