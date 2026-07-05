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

	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// runGraph coordinates the common reflector -> queue -> controller lifecycle.
//
// Callers validate graph shape and nil contexts before reaching this helper.
// The helper owns the shared shutdown policy: once either side exits, no more
// producer work should be accepted, the sibling side should be canceled, and
// the controller should be allowed to drain already queued items.
//
// The helper deliberately runs only one reflector/controller pair. It is not a
// higher-level supervisor: it does not restart components, classify domain
// errors, or recover panics.
func runGraph(
	ctx context.Context,
	queue *objectworkqueue.Queue,
	reflector *objectreflector.Reflector,
	controller *objectcontroller.Controller,
) error {
	return runComponents(ctx, queue, controller, reflector)
}

// runComponents coordinates one controller with one or more producer
// reflectors.
//
// Any reflector exit ends the graph because the controller no longer has a
// complete input set. The shared queue is shut down exactly once, then the
// controller drains already queued work before returning.
func runComponents(
	ctx context.Context,
	queue *objectworkqueue.Queue,
	controller *objectcontroller.Controller,
	reflectors ...*objectreflector.Reflector,
) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var shutdown sync.Once
	// stopGraph is shared by both result paths so queue shutdown is performed
	// exactly once even if both components return at nearly the same time.
	stopGraph := func() {
		shutdown.Do(func() {
			queue.ShutDown()
			cancel()
		})
	}

	componentCount := 1 + len(reflectors)
	results := make(chan graphRunResult, componentCount)
	for _, reflector := range reflectors {
		reflector := reflector
		go func() {
			results <- graphRunResult{err: reflector.Run(runCtx)}
		}()
	}
	go func() {
		results <- graphRunResult{err: controller.Run(runCtx)}
	}()

	first := <-results
	firstFatal := fatalBeforeRunnerCancel(ctx, first.err)
	stopGraph()

	fatals := []error{firstFatal}
	for i := 1; i < componentCount; i++ {
		result := <-results
		fatals = append(fatals, fatalAfterRunnerCancel(ctx, result.err))
	}

	return graphRunError(ctx, fatals...)
}

// graphRunResult keeps component completion reports uniform while preserving
// each component's error unchanged until classification.
type graphRunResult struct {
	err error
}

// fatalBeforeRunnerCancel classifies the first component result before the
// runner has canceled the sibling component. A context error is benign in this
// phase only when it comes from the parent context, because runner-initiated
// cancellation has not happened yet.
func fatalBeforeRunnerCancel(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if parentErr := ctx.Err(); parentErr != nil && errors.Is(err, parentErr) {
		return nil
	}

	return err
}

// fatalAfterRunnerCancel classifies the second component result after the
// runner has already canceled the shared child context. Plain cancellation from
// the sibling shutdown path is expected here; any non-context error still
// remains fatal.
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

// graphRunError turns classified component failures into the public runner
// result. Parent cancellation is reported only when no component produced a
// more specific fatal error.
func graphRunError(ctx context.Context, fatals ...error) error {
	err := errors.Join(fatals...)
	if err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
