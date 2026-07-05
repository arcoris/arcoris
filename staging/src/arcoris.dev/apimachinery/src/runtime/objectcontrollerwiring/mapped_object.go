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
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectreflectorsink"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// MappedObject is an assembled mapped-object controller runtime graph.
//
// MappedObject is passive like SameObject. It exposes the assembled components
// but does not start them, shut them down, or mutate mapping policy after
// construction.
type MappedObject struct {
	cache      *objectcache.Cache
	queue      *objectworkqueue.Queue
	reflector  *objectreflector.Reflector
	controller *objectcontroller.Controller
}

// NewMappedObject assembles a mapped-object controller graph.
//
// The assembled graph is:
//
//	objectreflector
//	  -> objectreflectorsink.Fanout(objectcache.Cache, objectenqueue.ReflectorSink)
//	  -> objectworkqueue.Queue
//	  -> objectcontroller.Controller
//
// Cache stores the watched source collection. Listed and Changed decide which
// mapped keys are enqueued. The fanout order stays cache first, enqueue second,
// so mapped reconciliations observe the reflected source state that produced
// their work items.
func NewMappedObject(config MappedObjectConfig) (*MappedObject, error) {
	cache, err := objectcache.New(config.Collection)
	if err != nil {
		return nil, err
	}

	queue, err := objectworkqueue.New(config.Queue)
	if err != nil {
		return nil, err
	}

	enqueueSink, err := objectenqueue.NewReflectorSink(objectenqueue.ReflectorSinkConfig{
		Queue:     queue,
		Predicate: config.Predicate,
		Listed:    config.Listed,
		Changed:   config.Changed,
	})
	if err != nil {
		return nil, err
	}

	fanout, err := objectreflectorsink.NewFanout(cache, enqueueSink)
	if err != nil {
		return nil, err
	}

	reflector, err := objectreflector.New(config.Source, config.Collection, fanout)
	if err != nil {
		return nil, err
	}

	controller, err := objectcontroller.New(config.Controller, queue, cache, config.Reconciler)
	if err != nil {
		return nil, err
	}

	return &MappedObject{
		cache:      cache,
		queue:      queue,
		reflector:  reflector,
		controller: controller,
	}, nil
}

// Cache returns the source-collection read model assembled for w.
func (w *MappedObject) Cache() *objectcache.Cache {
	if w == nil {
		return nil
	}

	return w.cache
}

// Queue returns the bounded mapped-work queue assembled for w.
func (w *MappedObject) Queue() *objectworkqueue.Queue {
	if w == nil {
		return nil
	}

	return w.queue
}

// Reflector returns the source-collection reflector assembled for w.
func (w *MappedObject) Reflector() *objectreflector.Reflector {
	if w == nil {
		return nil
	}

	return w.reflector
}

// Controller returns the worker lifecycle controller assembled for w.
func (w *MappedObject) Controller() *objectcontroller.Controller {
	if w == nil {
		return nil
	}

	return w.controller
}
