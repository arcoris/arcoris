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
)

func TestRunMappedObjectPanicsOnNilContext(t *testing.T) {
	graph := newMappedTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})

	defer func() {
		if recover() == nil {
			t.Fatal("RunMappedObject did not panic")
		}
	}()

	_ = RunMappedObject(nil, graph)
}

func TestRunMappedObjectRejectsInvalidGraph(t *testing.T) {
	tests := []struct {
		name  string
		graph func(t testing.TB) *MappedObject
	}{
		{
			name: "nil graph",
			graph: func(testing.TB) *MappedObject {
				return nil
			},
		},
		{
			name: "nil queue",
			graph: func(t testing.TB) *MappedObject {
				graph := newMappedTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})
				graph.queue = nil
				return graph
			},
		},
		{
			name: "nil reflector",
			graph: func(t testing.TB) *MappedObject {
				graph := newMappedTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})
				graph.reflector = nil
				return graph
			},
		},
		{
			name: "nil controller",
			graph: func(t testing.TB) *MappedObject {
				graph := newMappedTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})
				graph.controller = nil
				return graph
			},
		},
	}

	for _, tt := range tests {
		err := RunMappedObject(context.Background(), tt.graph(t))

		if !errors.Is(err, ErrInvalidMappedObject) {
			t.Fatalf("%s: error = %v; want errors.Is(%v)", tt.name, err, ErrInvalidMappedObject)
		}
	}
}

func TestRunMappedObjectShutsDownQueueWhenReflectorReturns(t *testing.T) {
	key := runTestKey("source-a")
	reflectorErr := errors.New("reflector failed")
	reconciler := &runTestReconciler{
		started: make(chan struct{}),
	}
	source := &runTestListerWatcher{
		read:            runTestRead(t, 10, runTestItem(key, 1, "source-a")),
		watchWaitStrict: reconciler.started,
		watchErr:        reflectorErr,
	}
	graph := newMappedTestGraph(t, source, reconciler)

	err := RunMappedObject(context.Background(), graph)

	requireErrorIs(t, err, reflectorErr)
	if !graph.Queue().IsShutDown() {
		t.Fatal("queue is not shut down")
	}
	requireReconcilerKeys(t, reconciler, key)
	requireQueueDrained(t, graph.Queue())
}
