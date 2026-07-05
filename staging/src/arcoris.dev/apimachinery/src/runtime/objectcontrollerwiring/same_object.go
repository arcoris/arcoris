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
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// SameObject is an assembled same-object controller runtime graph.
//
// SameObject is deliberately passive. It exposes the components created by
// NewSameObject, but does not start them, shut them down, or mutate wiring after
// construction.
type SameObject struct {
	// cache stores the reflected collection and serves controller snapshots.
	cache *objectcache.Cache
	// queue receives same-object work emitted by the enqueue sink.
	queue *objectworkqueue.Queue
	// indexes observe reflected state after cache and before enqueue.
	indexes []*objectindex.Index
	// reflector drives source reads and changes into cache and enqueue sinks.
	reflector *objectreflector.Reflector
	// controller consumes queue items and invokes the configured reconciler.
	controller *objectcontroller.Controller
}

// NewSameObject assembles the standard same-object controller graph.
//
// The assembled graph is:
//
//	objectreflector
//	  -> objectreflectorsink.Fanout(objectcache.Cache, indexes..., objectenqueue.ReflectorSink)
//	  -> objectworkqueue.Queue
//	  -> objectcontroller.Controller
//
// The fanout order is part of the contract. Cache is installed first, optional
// indexes second, and enqueue last so a controller worker and any mapper-owned
// index lookups observe the reflected state before work becomes visible.
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

	fanout, indexes, err := newInputFanout(cache, config.Indexes, enqueueSink)
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
		indexes:    indexes,
		reflector:  reflector,
		controller: controller,
	}, nil
}

// Cache returns the materialized read model assembled for w.
//
// A nil receiver returns nil so callers can inspect optional wiring values
// without defensive nil checks around the accessor itself.
func (w *SameObject) Cache() *objectcache.Cache {
	if w == nil {
		return nil
	}

	return w.cache
}

// Indexes returns the optional secondary indexes installed for w.
//
// A nil receiver returns nil. The returned slice is detached; index pointers
// are intentionally shared because callers close over the same instances for
// direct Lookup calls.
func (w *SameObject) Indexes() []*objectindex.Index {
	if w == nil {
		return nil
	}

	return append([]*objectindex.Index(nil), w.indexes...)
}

// Queue returns the bounded work queue assembled for w.
//
// The returned queue is the same queue used by the enqueue sink and controller.
// SameObject does not own shutdown policy; RunSameObject coordinates that when
// requested.
func (w *SameObject) Queue() *objectworkqueue.Queue {
	if w == nil {
		return nil
	}

	return w.queue
}

// Reflector returns the collection reflector assembled for w.
//
// The reflector is configured with fanout ordered as cache, then indexes, then
// enqueue.
func (w *SameObject) Reflector() *objectreflector.Reflector {
	if w == nil {
		return nil
	}

	return w.reflector
}

// Controller returns the worker lifecycle controller assembled for w.
//
// The controller consumes Queue and reconciles against Cache as its snapshot
// source.
func (w *SameObject) Controller() *objectcontroller.Controller {
	if w == nil {
		return nil
	}

	return w.controller
}
