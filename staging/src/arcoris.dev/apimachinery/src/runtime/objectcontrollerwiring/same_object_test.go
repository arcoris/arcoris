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

package objectcontrollerwiring_test

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectcontrollerwiring"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestNewSameObjectAcceptsValidConfig(t *testing.T) {
	wiring, err := objectcontrollerwiring.NewSameObject(validSameObjectConfig())
	requireNoError(t, err)

	if wiring == nil {
		t.Fatal("wiring is nil")
	}
	if wiring.Cache() == nil {
		t.Fatal("Cache() is nil")
	}
	if wiring.Queue() == nil {
		t.Fatal("Queue() is nil")
	}
	if wiring.Reflector() == nil {
		t.Fatal("Reflector() is nil")
	}
	if wiring.Controller() == nil {
		t.Fatal("Controller() is nil")
	}
}

func TestNewSameObjectRejectsInvalidConfigThroughDownstreamErrors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*objectcontrollerwiring.SameObjectConfig)
		target error
	}{
		{
			name: "nil source",
			mutate: func(config *objectcontrollerwiring.SameObjectConfig) {
				config.Source = nil
			},
			target: objectreflector.ErrNilSource,
		},
		{
			name: "invalid collection",
			mutate: func(config *objectcontrollerwiring.SameObjectConfig) {
				config.Collection = objectstore.ListRequest{}
			},
			target: objectstore.ErrInvalidListRequest,
		},
		{
			name: "nil reconciler",
			mutate: func(config *objectcontrollerwiring.SameObjectConfig) {
				config.Reconciler = nil
			},
			target: objectcontroller.ErrNilReconciler,
		},
		{
			name: "invalid queue options",
			mutate: func(config *objectcontrollerwiring.SameObjectConfig) {
				config.Queue = objectworkqueue.Options{}
			},
			target: objectworkqueue.ErrInvalidCapacity,
		},
		{
			name: "invalid controller options",
			mutate: func(config *objectcontrollerwiring.SameObjectConfig) {
				config.Controller = objectcontroller.Options{}
			},
			target: objectcontroller.ErrInvalidWorkers,
		},
	}

	for _, tt := range tests {
		config := validSameObjectConfig()
		tt.mutate(&config)

		_, err := objectcontrollerwiring.NewSameObject(config)

		if !errors.Is(err, tt.target) {
			t.Fatalf("%s: error = %v; want errors.Is(%v)", tt.name, err, tt.target)
		}
	}
}

func TestSameObjectGettersReturnNilForNilReceiver(t *testing.T) {
	var wiring *objectcontrollerwiring.SameObject

	if wiring.Cache() != nil {
		t.Fatal("Cache() is not nil")
	}
	if wiring.Queue() != nil {
		t.Fatal("Queue() is not nil")
	}
	if wiring.Reflector() != nil {
		t.Fatal("Reflector() is not nil")
	}
	if wiring.Controller() != nil {
		t.Fatal("Controller() is not nil")
	}
}

func validSameObjectConfig() objectcontrollerwiring.SameObjectConfig {
	return objectcontrollerwiring.SameObjectConfig{
		Source:     &scriptedListerWatcher{},
		Collection: wiringCollection(),
		Reconciler: &recordingReconciler{},
		Queue: objectworkqueue.Options{
			Capacity: 4,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
	}
}

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
}
