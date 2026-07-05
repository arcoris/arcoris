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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestNewMappedObjectAcceptsValidConfig(t *testing.T) {
	graph, err := NewMappedObject(validMappedObjectConfig())
	requireNoError(t, err)

	if graph == nil {
		t.Fatal("graph is nil")
	}
	if graph.Cache() == nil {
		t.Fatal("Cache() is nil")
	}
	if graph.Queue() == nil {
		t.Fatal("Queue() is nil")
	}
	if graph.Reflector() == nil {
		t.Fatal("Reflector() is nil")
	}
	if graph.Controller() == nil {
		t.Fatal("Controller() is nil")
	}
}

func TestNewMappedObjectAcceptsIndexes(t *testing.T) {
	index := newDesiredObjectIndex(t)
	config := validMappedObjectConfig()
	config.Indexes = []*objectindex.Index{index}

	graph, err := NewMappedObject(config)
	requireNoError(t, err)

	indexes := graph.Indexes()
	if len(indexes) != 1 || indexes[0] != index {
		t.Fatalf("Indexes() = %#v; want configured index", indexes)
	}
}

func TestMappedObjectIndexesReturnsDetachedSlice(t *testing.T) {
	index := newDesiredObjectIndex(t)
	config := validMappedObjectConfig()
	config.Indexes = []*objectindex.Index{index}
	graph, err := NewMappedObject(config)
	requireNoError(t, err)

	indexes := graph.Indexes()
	indexes[0] = nil

	if graph.Indexes()[0] == nil {
		t.Fatal("Indexes() exposed internal slice")
	}
}

func TestNewMappedObjectRejectsInvalidConfigThroughDownstreamErrors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*MappedObjectConfig)
		target error
	}{
		{
			name: "nil source",
			mutate: func(config *MappedObjectConfig) {
				config.Source = nil
			},
			target: objectreflector.ErrNilSource,
		},
		{
			name: "invalid collection",
			mutate: func(config *MappedObjectConfig) {
				config.Collection = objectstore.ListRequest{}
			},
			target: objectstore.ErrInvalidListRequest,
		},
		{
			name: "nil reconciler",
			mutate: func(config *MappedObjectConfig) {
				config.Reconciler = nil
			},
			target: objectcontroller.ErrNilReconciler,
		},
		{
			name: "nil listed mapper",
			mutate: func(config *MappedObjectConfig) {
				config.Listed = nil
			},
			target: objectenqueue.ErrNilListItemMapper,
		},
		{
			name: "nil changed mapper",
			mutate: func(config *MappedObjectConfig) {
				config.Changed = nil
			},
			target: objectenqueue.ErrNilMapper,
		},
		{
			name: "invalid queue options",
			mutate: func(config *MappedObjectConfig) {
				config.Queue = objectworkqueue.Options{}
			},
			target: objectworkqueue.ErrInvalidCapacity,
		},
		{
			name: "invalid controller options",
			mutate: func(config *MappedObjectConfig) {
				config.Controller = objectcontroller.Options{}
			},
			target: objectcontroller.ErrInvalidWorkers,
		},
		{
			name: "nil index",
			mutate: func(config *MappedObjectConfig) {
				config.Indexes = []*objectindex.Index{nil}
			},
			target: ErrNilIndex,
		},
	}

	for _, tt := range tests {
		config := validMappedObjectConfig()
		tt.mutate(&config)

		_, err := NewMappedObject(config)

		if !errors.Is(err, tt.target) {
			t.Fatalf("%s: error = %v; want errors.Is(%v)", tt.name, err, tt.target)
		}
	}
}

func TestMappedObjectGettersReturnNilForNilReceiver(t *testing.T) {
	var graph *MappedObject

	if graph.Cache() != nil {
		t.Fatal("Cache() is not nil")
	}
	if graph.Queue() != nil {
		t.Fatal("Queue() is not nil")
	}
	if graph.Reflector() != nil {
		t.Fatal("Reflector() is not nil")
	}
	if graph.Controller() != nil {
		t.Fatal("Controller() is not nil")
	}
	if graph.Indexes() != nil {
		t.Fatal("Indexes() is not nil")
	}
}

func newMappedTestGraph(
	t testing.TB,
	source *runTestListerWatcher,
	reconciler *runTestReconciler,
) *MappedObject {
	t.Helper()

	graph, err := NewMappedObject(MappedObjectConfig{
		Source:     source,
		Collection: runTestCollection(),
		Reconciler: reconciler,
		Queue: objectworkqueue.Options{
			Capacity: 8,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
		Listed:  objectenqueue.ListedObject(),
		Changed: objectenqueue.ChangedObject(),
	})
	requireNoError(t, err)

	return graph
}

// validMappedObjectConfig keeps constructor tests focused on one invalid field
// at a time. The mappers intentionally emit no work because constructor tests
// care only about graph assembly and downstream validation.
func validMappedObjectConfig() MappedObjectConfig {
	return MappedObjectConfig{
		Source:     &runTestListerWatcher{},
		Collection: runTestCollection(),
		Reconciler: &runTestReconciler{},
		Queue: objectworkqueue.Options{
			Capacity: 4,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
		Listed: objectenqueue.ListItemMapperFunc(
			func(objectstore.ListItem, objectenqueue.EmitFunc) error {
				return nil
			},
		),
		Changed: objectenqueue.MapperFunc(
			func(objectstore.Change, objectenqueue.EmitFunc) error {
				return nil
			},
		),
	}
}
