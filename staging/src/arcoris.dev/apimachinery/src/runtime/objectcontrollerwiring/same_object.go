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

// SameObject is an assembled same-object controller runtime graph.
type SameObject struct {
	cache      *objectcache.Cache
	queue      *objectworkqueue.Queue
	reflector  *objectreflector.Reflector
	controller *objectcontroller.Controller
}

// NewSameObject assembles the standard same-object controller graph.
//
// The reflector fanout is deliberately ordered as cache first and enqueue sink
// second. This preserves the invariant that a controller worker can observe the
// reflected state in its cache snapshot when it receives the corresponding
// object work item.
func NewSameObject(config SameObjectConfig) (*SameObject, error) {
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
		Listed:    objectenqueue.ListedObject(),
		Changed:   objectenqueue.ChangedObject(),
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

	return &SameObject{
		cache:      cache,
		queue:      queue,
		reflector:  reflector,
		controller: controller,
	}, nil
}

// Cache returns the materialized read model assembled for w.
func (w *SameObject) Cache() *objectcache.Cache {
	if w == nil {
		return nil
	}

	return w.cache
}

// Queue returns the bounded work queue assembled for w.
func (w *SameObject) Queue() *objectworkqueue.Queue {
	if w == nil {
		return nil
	}

	return w.queue
}

// Reflector returns the collection reflector assembled for w.
func (w *SameObject) Reflector() *objectreflector.Reflector {
	if w == nil {
		return nil
	}

	return w.reflector
}

// Controller returns the worker lifecycle controller assembled for w.
func (w *SameObject) Controller() *objectcontroller.Controller {
	if w == nil {
		return nil
	}

	return w.controller
}
