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
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestNewMultiSourceAcceptsPrimaryOnlyConfig(t *testing.T) {
	graph, err := NewMultiSource(validMultiSourceConfig())
	requireNoError(t, err)

	requireValidMultiSourceGraph(t, graph)
	if got := len(graph.Secondary()); got != 0 {
		t.Fatalf("secondary inputs = %d; want 0", got)
	}
	if got := len(graph.Inputs()); got != 1 {
		t.Fatalf("inputs = %d; want 1", got)
	}
}

func TestNewMultiSourceAcceptsPrimaryPlusSecondaryConfigs(t *testing.T) {
	config := validMultiSourceConfig()
	config.Secondary = []InputConfig{
		validInputConfig(&runTestListerWatcher{}),
		validInputConfig(&runTestListerWatcher{}),
	}

	graph, err := NewMultiSource(config)
	requireNoError(t, err)

	requireValidMultiSourceGraph(t, graph)
	if got := len(graph.Secondary()); got != 2 {
		t.Fatalf("secondary inputs = %d; want 2", got)
	}
	if got := len(graph.Inputs()); got != 3 {
		t.Fatalf("inputs = %d; want 3", got)
	}
}

func TestMultiSourceExposesDetachedInputSlices(t *testing.T) {
	config := validMultiSourceConfig()
	config.Secondary = []InputConfig{validInputConfig(&runTestListerWatcher{})}
	graph, err := NewMultiSource(config)
	requireNoError(t, err)

	secondary := graph.Secondary()
	secondary[0] = nil
	if graph.Secondary()[0] == nil {
		t.Fatal("Secondary() exposed internal slice")
	}

	inputs := graph.Inputs()
	inputs[0] = nil
	if graph.Inputs()[0] == nil {
		t.Fatal("Inputs() exposed internal slice")
	}
}

func TestNewMultiSourceRejectsInvalidConfigThroughDownstreamErrors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*MultiSourceConfig)
		target error
	}{
		{
			name: "nil primary source",
			mutate: func(config *MultiSourceConfig) {
				config.Primary.Source = nil
			},
			target: objectreflector.ErrNilSource,
		},
		{
			name: "invalid primary collection",
			mutate: func(config *MultiSourceConfig) {
				config.Primary.Collection = objectstore.ListRequest{}
			},
			target: objectstore.ErrInvalidListRequest,
		},
		{
			name: "nil primary listed mapper",
			mutate: func(config *MultiSourceConfig) {
				config.Primary.Listed = nil
			},
			target: objectenqueue.ErrNilListItemMapper,
		},
		{
			name: "nil primary changed mapper",
			mutate: func(config *MultiSourceConfig) {
				config.Primary.Changed = nil
			},
			target: objectenqueue.ErrNilMapper,
		},
		{
			name: "nil secondary source",
			mutate: func(config *MultiSourceConfig) {
				secondary := validInputConfig(&runTestListerWatcher{})
				secondary.Source = nil
				config.Secondary = []InputConfig{secondary}
			},
			target: objectreflector.ErrNilSource,
		},
		{
			name: "invalid secondary collection",
			mutate: func(config *MultiSourceConfig) {
				secondary := validInputConfig(&runTestListerWatcher{})
				secondary.Collection = objectstore.ListRequest{}
				config.Secondary = []InputConfig{secondary}
			},
			target: objectstore.ErrInvalidListRequest,
		},
		{
			name: "nil secondary listed mapper",
			mutate: func(config *MultiSourceConfig) {
				secondary := validInputConfig(&runTestListerWatcher{})
				secondary.Listed = nil
				config.Secondary = []InputConfig{secondary}
			},
			target: objectenqueue.ErrNilListItemMapper,
		},
		{
			name: "nil secondary changed mapper",
			mutate: func(config *MultiSourceConfig) {
				secondary := validInputConfig(&runTestListerWatcher{})
				secondary.Changed = nil
				config.Secondary = []InputConfig{secondary}
			},
			target: objectenqueue.ErrNilMapper,
		},
		{
			name: "nil reconciler",
			mutate: func(config *MultiSourceConfig) {
				config.Reconciler = nil
			},
			target: objectcontroller.ErrNilReconciler,
		},
		{
			name: "invalid queue options",
			mutate: func(config *MultiSourceConfig) {
				config.Queue = objectworkqueue.Options{}
			},
			target: objectworkqueue.ErrInvalidCapacity,
		},
		{
			name: "invalid controller options",
			mutate: func(config *MultiSourceConfig) {
				config.Controller = objectcontroller.Options{}
			},
			target: objectcontroller.ErrInvalidWorkers,
		},
	}

	for _, tt := range tests {
		config := validMultiSourceConfig()
		tt.mutate(&config)

		_, err := NewMultiSource(config)

		if !errors.Is(err, tt.target) {
			t.Fatalf("%s: error = %v; want errors.Is(%v)", tt.name, err, tt.target)
		}
	}
}

func TestMultiSourceGettersReturnNilForNilReceiver(t *testing.T) {
	var graph *MultiSource

	if graph.Primary() != nil {
		t.Fatal("Primary() is not nil")
	}
	if graph.Secondary() != nil {
		t.Fatal("Secondary() is not nil")
	}
	if graph.Inputs() != nil {
		t.Fatal("Inputs() is not nil")
	}
	if graph.Queue() != nil {
		t.Fatal("Queue() is not nil")
	}
	if graph.Controller() != nil {
		t.Fatal("Controller() is not nil")
	}
}

func requireValidMultiSourceGraph(t testing.TB, graph *MultiSource) {
	t.Helper()

	if graph == nil {
		t.Fatal("graph is nil")
	}
	if graph.Primary() == nil {
		t.Fatal("Primary() is nil")
	}
	if graph.Primary().Cache() == nil {
		t.Fatal("Primary().Cache() is nil")
	}
	if graph.Primary().Reflector() == nil {
		t.Fatal("Primary().Reflector() is nil")
	}
	if graph.Queue() == nil {
		t.Fatal("Queue() is nil")
	}
	if graph.Controller() == nil {
		t.Fatal("Controller() is nil")
	}
}

// validMultiSourceConfig is a minimal primary-only graph config. Individual
// tests mutate one field at a time so the expected downstream error remains
// obvious.
func validMultiSourceConfig() MultiSourceConfig {
	return MultiSourceConfig{
		Primary:    validInputConfig(&runTestListerWatcher{}),
		Reconciler: &runTestReconciler{},
		Queue: objectworkqueue.Options{
			Capacity: 8,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
	}
}

func validInputConfig(source *runTestListerWatcher) InputConfig {
	return InputConfig{
		Source:     source,
		Collection: runTestCollection(),
		Listed:     objectenqueue.ListedObject(),
		Changed:    objectenqueue.ChangedObject(),
	}
}
