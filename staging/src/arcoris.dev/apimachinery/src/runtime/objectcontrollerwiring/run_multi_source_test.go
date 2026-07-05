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
	"testing"
	"time"
)

func TestRunMultiSourcePanicsOnNilContext(t *testing.T) {
	graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, nil, &runTestReconciler{})

	defer func() {
		if recover() == nil {
			t.Fatal("RunMultiSource did not panic")
		}
	}()

	_ = RunMultiSource(nil, graph)
}

func TestRunMultiSourceRejectsInvalidGraph(t *testing.T) {
	tests := []struct {
		name  string
		graph func(t testing.TB) *MultiSource
	}{
		{
			name: "nil graph",
			graph: func(testing.TB) *MultiSource {
				return nil
			},
		},
		{
			name: "nil queue",
			graph: func(t testing.TB) *MultiSource {
				graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, nil, &runTestReconciler{})
				graph.queue = nil
				return graph
			},
		},
		{
			name: "nil controller",
			graph: func(t testing.TB) *MultiSource {
				graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, nil, &runTestReconciler{})
				graph.controller = nil
				return graph
			},
		},
		{
			name: "nil primary",
			graph: func(t testing.TB) *MultiSource {
				graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, nil, &runTestReconciler{})
				graph.primary = nil
				return graph
			},
		},
		{
			name: "nil primary reflector",
			graph: func(t testing.TB) *MultiSource {
				graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, nil, &runTestReconciler{})
				graph.primary.reflector = nil
				return graph
			},
		},
		{
			name: "nil secondary",
			graph: func(t testing.TB) *MultiSource {
				graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, []*runTestListerWatcher{{}}, &runTestReconciler{})
				graph.secondary[0] = nil
				return graph
			},
		},
		{
			name: "nil secondary reflector",
			graph: func(t testing.TB) *MultiSource {
				graph := newMultiSourceRunTestGraph(t, &runTestListerWatcher{}, []*runTestListerWatcher{{}}, &runTestReconciler{})
				graph.secondary[0].reflector = nil
				return graph
			},
		},
	}

	for _, tt := range tests {
		err := RunMultiSource(context.Background(), tt.graph(t))

		if !errors.Is(err, ErrInvalidMultiSource) {
			t.Fatalf("%s: error = %v; want errors.Is(%v)", tt.name, err, ErrInvalidMultiSource)
		}
	}
}

func TestRunMultiSourceReturnsParentContextError(t *testing.T) {
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 10),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	secondary := &runTestListerWatcher{
		read:        runTestRead(t, 10),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newMultiSourceRunTestGraph(t, primary, []*runTestListerWatcher{secondary}, &runTestReconciler{})
	ctx, cancel := context.WithCancel(context.Background())

	result := runMultiSourceAsync(ctx, graph)
	waitForSignal(t, primary.watchCalled)
	waitForSignal(t, secondary.watchCalled)
	cancel()

	requireMultiSourceRunResult(t, result, context.Canceled)
}

func TestRunMultiSourceReturnsSecondarySourceErrorAndStopsGraph(t *testing.T) {
	sourceErr := errors.New("secondary source failed")
	primaryStream := runTestWaitingStream()
	primaryWatchOpened := make(chan struct{})
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 10),
		stream:      primaryStream,
		watchCalled: primaryWatchOpened,
	}
	secondary := &runTestListerWatcher{
		read:            runTestRead(t, 10),
		watchWaitStrict: primaryWatchOpened,
		watchErr:        sourceErr,
	}
	graph := newMultiSourceRunTestGraph(t, primary, []*runTestListerWatcher{secondary}, &runTestReconciler{})

	err := RunMultiSource(context.Background(), graph)

	requireErrorIs(t, err, sourceErr)
	if !graph.Queue().IsShutDown() {
		t.Fatal("queue is not shut down")
	}
	waitForSignal(t, primaryStream.done)
}

func newMultiSourceRunTestGraph(
	t testing.TB,
	primary *runTestListerWatcher,
	secondary []*runTestListerWatcher,
	reconciler *runTestReconciler,
) *MultiSource {
	t.Helper()

	config := validMultiSourceConfig()
	config.Primary = validInputConfig(primary)
	config.Reconciler = reconciler
	for _, source := range secondary {
		config.Secondary = append(config.Secondary, validInputConfig(source))
	}

	graph, err := NewMultiSource(config)
	requireNoError(t, err)

	return graph
}

type multiSourceRunResult struct {
	err error
}

// runMultiSourceAsync starts the real multi-source runner for tests that need
// to coordinate cancellation after all input reflectors are waiting.
func runMultiSourceAsync(ctx context.Context, graph *MultiSource) <-chan multiSourceRunResult {
	result := make(chan multiSourceRunResult, 1)
	go func() {
		result <- multiSourceRunResult{err: RunMultiSource(ctx, graph)}
	}()

	return result
}

func requireMultiSourceRunResult(t testing.TB, result <-chan multiSourceRunResult, target error) {
	t.Helper()

	requireErrorIs(t, readMultiSourceRunResult(t, result), target)
}

func readMultiSourceRunResult(t testing.TB, result <-chan multiSourceRunResult) error {
	t.Helper()

	select {
	case result := <-result:
		return result.err
	case <-time.After(5 * time.Second):
		t.Fatal("RunMultiSource did not return")
		return nil
	}
}
