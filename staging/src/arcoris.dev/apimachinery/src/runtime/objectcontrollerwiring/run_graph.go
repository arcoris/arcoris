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
func runGraph(
	ctx context.Context,
	queue *objectworkqueue.Queue,
	reflector *objectreflector.Reflector,
	controller *objectcontroller.Controller,
) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var shutdown sync.Once
	stopGraph := func() {
		shutdown.Do(func() {
			queue.ShutDown()
			cancel()
		})
	}

	results := make(chan graphRunResult, 2)
	go func() {
		results <- graphRunResult{err: reflector.Run(runCtx)}
	}()
	go func() {
		results <- graphRunResult{err: controller.Run(runCtx)}
	}()

	first := <-results
	firstFatal := fatalBeforeRunnerCancel(ctx, first.err)
	stopGraph()

	second := <-results
	secondFatal := fatalAfterRunnerCancel(ctx, second.err)

	return graphRunError(ctx, firstFatal, secondFatal)
}

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

func graphRunError(ctx context.Context, first error, second error) error {
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
