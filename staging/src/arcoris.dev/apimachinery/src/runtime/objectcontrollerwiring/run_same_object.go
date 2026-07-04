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

package objectcontrollerwiring

import (
	"context"
	"errors"
	"sync"
)

// RunSameObject runs one already assembled same-object controller graph.
//
// RunSameObject starts the reflector and controller with a shared child
// context. When either side exits, it shuts down the queue, cancels the sibling
// side, waits for both sides to stop, and returns the fatal error, joined fatal
// errors, ctx.Err, or nil according to the observed shutdown cause.
//
// RunSameObject panics when ctx is nil. It does not recover panics from lower
// components.
func RunSameObject(ctx context.Context, graph *SameObject) error {
	if ctx == nil {
		panic("objectcontrollerwiring: nil context")
	}
	if err := validateRunnableSameObject(graph); err != nil {
		return err
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var shutdown sync.Once
	stopGraph := func() {
		shutdown.Do(func() {
			graph.Queue().ShutDown()
			cancel()
		})
	}

	results := make(chan sameObjectRunResult, 2)
	go func() {
		results <- sameObjectRunResult{err: graph.Reflector().Run(runCtx)}
	}()
	go func() {
		results <- sameObjectRunResult{err: graph.Controller().Run(runCtx)}
	}()

	first := <-results
	firstFatal := fatalBeforeRunnerCancel(ctx, first.err)
	stopGraph()

	second := <-results
	secondFatal := fatalAfterRunnerCancel(ctx, second.err)

	return sameObjectRunError(ctx, firstFatal, secondFatal)
}

type sameObjectRunResult struct {
	err error
}

// validateRunnableSameObject checks only the components coordinated by the
// runner. Cache is intentionally not required here because the controller
// already owns its snapshot source dependency.
func validateRunnableSameObject(graph *SameObject) error {
	if graph == nil ||
		graph.Queue() == nil ||
		graph.Reflector() == nil ||
		graph.Controller() == nil {
		return ErrInvalidSameObject
	}

	return nil
}

// fatalBeforeRunnerCancel classifies the first component result. At this point
// the runner has not canceled the child context yet, so a context error is
// benign only when the parent context is already canceled.
func fatalBeforeRunnerCancel(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if parentErr := ctx.Err(); parentErr != nil && errors.Is(err, parentErr) {
		return nil
	}

	return err
}

// fatalAfterRunnerCancel classifies the second component result. The runner has
// canceled the child context by this point, so ordinary cancellation from the
// sibling shutdown path is benign.
func fatalAfterRunnerCancel(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if parentErr := ctx.Err(); parentErr != nil && errors.Is(err, parentErr) {
		return nil
	}

	return err
}

func sameObjectRunError(ctx context.Context, first error, second error) error {
	switch {
	case first != nil && second != nil:
		return errors.Join(first, second)
	case first != nil:
		return first
	case second != nil:
		return second
	case ctx.Err() != nil:
		return ctx.Err()
	default:
		return nil
	}
}
