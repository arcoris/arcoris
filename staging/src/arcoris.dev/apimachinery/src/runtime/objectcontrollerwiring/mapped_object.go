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

// MappedObject is an assembled mapped-object controller runtime graph.
//
// MappedObject is passive like SameObject. It records the constructed
// components and mapping policy at assembly time, then leaves lifecycle
// coordination to RunMappedObject or to a caller that wants to drive the
// components manually.
type MappedObject struct {
	// cache stores the watched source collection, not the mapped target
	// collection. Mapped reconcilers receive this cache as their snapshot source.
	cache *objectcache.Cache
	// queue receives work items emitted by the caller-provided source mappers.
	queue *objectworkqueue.Queue
	// indexes observe source state before caller-provided mappers emit work.
	indexes []*objectindex.Index
	// reflector drives the source list-watch stream into cache and enqueue sinks.
	reflector *objectreflector.Reflector
	// controller consumes mapped queue items and invokes the configured reconciler.
	controller *objectcontroller.Controller
}

// NewMappedObject assembles a mapped-object controller graph.
//
// The assembled graph is:
//
//	objectreflector
//	  -> objectreflectorsink.Fanout(objectcache.Cache, indexes..., objectenqueue.ReflectorSink)
//	  -> objectworkqueue.Queue
//	  -> objectcontroller.Controller
//
// Cache stores the watched source collection. Listed and Changed decide which
// mapped keys are enqueued. The fanout order stays cache first, indexes second,
// enqueue last, so mappers can use indexes that already include the reflected
// state that produced their work items.
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

	return &MappedObject{
		cache:      cache,
		queue:      queue,
		indexes:    indexes,
		reflector:  reflector,
		controller: controller,
	}, nil
}

// Cache returns the source-collection read model assembled for w.
//
// A nil receiver returns nil, matching the other graph accessors. The returned
// cache is intentionally the watched source cache; target object reads, when a
// mapped reconciler needs them, must be provided through explicit reconciler
// dependencies outside this wiring value.
func (w *MappedObject) Cache() *objectcache.Cache {
	if w == nil {
		return nil
	}

	return w.cache
}

// Indexes returns the optional secondary indexes installed for w.
//
// A nil receiver returns nil. The returned slice is detached; the index
// instances are shared with caller-owned mapper closures by design.
func (w *MappedObject) Indexes() []*objectindex.Index {
	if w == nil {
		return nil
	}

	return append([]*objectindex.Index(nil), w.indexes...)
}

// Queue returns the bounded mapped-work queue assembled for w.
//
// The returned queue is shared by the enqueue sink and controller. MappedObject
// does not shut the queue down by itself; RunMappedObject coordinates that
// boundary when the graph is run through this package.
func (w *MappedObject) Queue() *objectworkqueue.Queue {
	if w == nil {
		return nil
	}

	return w.queue
}

// Reflector returns the source-collection reflector assembled for w.
//
// The reflector is configured with fanout ordered as cache, indexes, then
// enqueue. That ordering lets mapped mappers and reconcilers observe reflected
// source state before mapped requests are processed from the queue.
func (w *MappedObject) Reflector() *objectreflector.Reflector {
	if w == nil {
		return nil
	}

	return w.reflector
}

// Controller returns the worker lifecycle controller assembled for w.
//
// The controller consumes Queue and reconciles mapped requests against Cache as
// its snapshot source.
func (w *MappedObject) Controller() *objectcontroller.Controller {
	if w == nil {
		return nil
	}

	return w.controller
}
