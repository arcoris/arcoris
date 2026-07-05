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
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectreflectorsink"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// Input is one assembled watched input inside a MultiSource graph.
//
// Input is intentionally passive. It exposes the input cache and reflector, but
// not the enqueue sink or fanout used internally to connect reflected state to
// the graph's shared queue.
type Input struct {
	// cache stores this input collection only.
	cache *objectcache.Cache
	// reflector drives this input source into cache and shared-queue enqueue.
	reflector *objectreflector.Reflector
}

// newInput assembles one watched input against an already-created shared queue.
//
// The input fanout is always cache first and enqueue sink second. That preserves
// the same visibility invariant as SameObject and MappedObject for every input:
// an emitted work item is visible only after this input's cache has observed the
// reflected state that produced it.
func newInput(config InputConfig, queue *objectworkqueue.Queue) (*Input, error) {
	cache, err := objectcache.New(config.Collection)
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

	return &Input{
		cache:     cache,
		reflector: reflector,
	}, nil
}

// Cache returns the read model assembled for this input.
//
// A nil receiver returns nil. In MultiSource, the primary input cache is the
// controller snapshot source; secondary input caches are available as explicit
// side dependencies for callers that need them.
func (i *Input) Cache() *objectcache.Cache {
	if i == nil {
		return nil
	}

	return i.cache
}

// Reflector returns the reflector assembled for this input.
//
// The reflector is configured with input-local fanout ordered as cache first
// and enqueue sink second, so work from this input is visible only after the
// input cache observes the reflected state.
func (i *Input) Reflector() *objectreflector.Reflector {
	if i == nil {
		return nil
	}

	return i.reflector
}
