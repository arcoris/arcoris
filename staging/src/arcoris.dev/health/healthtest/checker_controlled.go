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

package healthtest

import (
	"context"
	"sync"

	"arcoris.dev/health"
)

// BlockingChecker blocks until released or until its context is canceled.
//
// BlockingChecker is useful for cancellation, timeout, and concurrency tests. It
// is safe for concurrent Check calls. Release unblocks all current and future
// calls with the same result.
type BlockingChecker struct {
	// name is the stable checker name returned by Name and used in cancellation
	// results.
	name string

	// mu protects calls and result because Check may run in several evaluator
	// workers while Release updates the result.
	mu    sync.Mutex
	calls int

	// started closes after the first Check call reaches the blocking point.
	started     chan struct{}
	startedOnce sync.Once

	// release closes when Release is called. Closing instead of sending wakes all
	// current and future Check calls.
	release     chan struct{}
	releaseOnce sync.Once

	// result is returned after release. It is protected by mu.
	result health.Result
}

// NewBlockingChecker returns a controlled blocking checker.
//
// The default release result is HealthyResult(name), so tests that only need to
// unblock the checker can call Release(HealthyResult(name)) explicitly or use a
// custom result for status-specific assertions.
func NewBlockingChecker(name string) *BlockingChecker {
	return &BlockingChecker{
		name:    name,
		started: make(chan struct{}),
		release: make(chan struct{}),
		result:  HealthyResult(name),
	}
}

// Name returns the configured checker name.
func (c *BlockingChecker) Name() string {
	if c == nil {
		return ""
	}

	return c.name
}

// Check blocks until Release is called or ctx is canceled.
//
// Context cancellation returns an UNKNOWN result with ReasonCanceled, matching
// the package-health convention for observations that cannot complete
// reliably. A nil context is treated as context.Background for convenience in
// tests that are not focused on cancellation.
func (c *BlockingChecker) Check(ctx context.Context) health.Result {
	if c == nil {
		return health.Unknown("", health.ReasonNotObserved, "nil blocking checker")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	c.mu.Lock()
	c.calls++
	c.mu.Unlock()
	c.startedOnce.Do(func() { close(c.started) })

	select {
	case <-c.release:
		c.mu.Lock()
		result := c.result
		c.mu.Unlock()
		return result
	case <-ctx.Done():
		return health.Unknown(c.name, health.ReasonCanceled, "health check canceled")
	}
}

// Release unblocks Check calls with result.
//
// Release is idempotent with respect to unblocking: the first call closes the
// release channel. Later calls may update result for future callers that have
// not yet read it, but tests should normally release once for deterministic
// assertions.
func (c *BlockingChecker) Release(result health.Result) {
	if c == nil {
		return
	}

	c.mu.Lock()
	c.result = result
	c.mu.Unlock()
	c.releaseOnce.Do(func() { close(c.release) })
}

// Started returns a channel closed after the first Check call starts.
//
// Tests can wait on Started instead of sleeping before asserting cancellation,
// timeout, or concurrent evaluator behavior.
func (c *BlockingChecker) Started() <-chan struct{} {
	if c == nil {
		closed := make(chan struct{})
		close(closed)
		return closed
	}

	return c.started
}

// Calls returns the number of Check calls.
func (c *BlockingChecker) Calls() int {
	if c == nil {
		return 0
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.calls
}

// SequenceChecker returns scripted results in order.
//
// The final result is sticky: once the sequence is exhausted, later calls keep
// returning the last result. The checker is safe for concurrent use.
type SequenceChecker struct {
	// name is the checker name returned by Name and used to fill unnamed results.
	name string

	// mu protects results and calls across concurrent evaluator workers.
	mu sync.Mutex

	// results is the remaining sticky sequence.
	results []health.Result

	// calls records Check invocations.
	calls int
}

// NewSequenceChecker returns a checker backed by results.
//
// The input slice is copied so later caller mutations cannot rewrite the
// sequence observed by evaluator tests.
func NewSequenceChecker(name string, results ...health.Result) *SequenceChecker {
	copied := make([]health.Result, len(results))
	copy(copied, results)

	return &SequenceChecker{name: name, results: copied}
}

// Name returns the configured checker name.
func (c *SequenceChecker) Name() string {
	if c == nil {
		return ""
	}

	return c.name
}

// Check returns the next scripted result.
//
// When a scripted result has an empty Name, Check fills it with the checker name
// so tests can use terse result literals without losing normal checker
// ownership semantics.
func (c *SequenceChecker) Check(context.Context) health.Result {
	if c == nil {
		return health.Unknown("", health.ReasonNotObserved, "nil sequence checker")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.calls++
	if len(c.results) == 0 {
		return UnknownResult(c.name, health.ReasonNotObserved)
	}

	result := c.results[0]
	if len(c.results) > 1 {
		c.results = c.results[1:]
	}
	if result.Name == "" {
		result.Name = c.name
	}

	return result
}

// Calls returns the number of Check calls.
func (c *SequenceChecker) Calls() int {
	if c == nil {
		return 0
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.calls
}
