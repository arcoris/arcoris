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
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// MultiSource is an assembled controller graph with multiple watched inputs.
//
// MultiSource has one shared queue and one controller. Every input has its own
// cache and reflector, and the primary input cache is the controller snapshot
// source.
type MultiSource struct {
	// primary provides the snapshot source used by controller.
	primary *Input
	// secondary contains optional additional inputs that emit into queue.
	secondary []*Input
	// queue is shared by every input and consumed by controller.
	queue *objectworkqueue.Queue
	// controller consumes queue and reconciles against primary.cache.
	controller *objectcontroller.Controller
}

// NewMultiSource assembles a multi-source controller graph.
//
// Construction creates exactly one shared queue and one controller. The primary
// input and every secondary input are assembled with that shared queue. The
// controller receives the primary input cache as its snapshot source; secondary
// caches remain separate and are exposed only through accessors.
func NewMultiSource(config MultiSourceConfig) (*MultiSource, error) {
	queue, err := objectworkqueue.New(config.Queue)
	if err != nil {
		return nil, err
	}

	primary, err := newInput(config.Primary, queue)
	if err != nil {
		return nil, err
	}

	secondary := make([]*Input, 0, len(config.Secondary))
	for _, inputConfig := range config.Secondary {
		input, err := newInput(inputConfig, queue)
		if err != nil {
			return nil, err
		}
		secondary = append(secondary, input)
	}

	controller, err := objectcontroller.New(config.Controller, queue, primary.Cache(), config.Reconciler)
	if err != nil {
		return nil, err
	}

	return &MultiSource{
		primary:    primary,
		secondary:  secondary,
		queue:      queue,
		controller: controller,
	}, nil
}

// Primary returns the primary watched input.
//
// A nil receiver returns nil. The primary input cache is the snapshot source
// used by Controller.
func (g *MultiSource) Primary() *Input {
	if g == nil {
		return nil
	}

	return g.primary
}

// Secondary returns the secondary watched inputs in construction order.
//
// A nil receiver returns nil. The returned slice is detached so callers cannot
// mutate the graph's internal input ordering.
func (g *MultiSource) Secondary() []*Input {
	if g == nil {
		return nil
	}

	return append([]*Input(nil), g.secondary...)
}

// Inputs returns all watched inputs in run order: primary first, then secondary.
//
// A nil receiver returns nil. The returned slice is detached.
func (g *MultiSource) Inputs() []*Input {
	if g == nil {
		return nil
	}

	inputs := make([]*Input, 0, 1+len(g.secondary))
	inputs = append(inputs, g.primary)
	inputs = append(inputs, g.secondary...)

	return inputs
}

// Queue returns the shared work queue assembled for g.
//
// A nil receiver returns nil. All inputs emit into this queue, and Controller
// consumes from it.
func (g *MultiSource) Queue() *objectworkqueue.Queue {
	if g == nil {
		return nil
	}

	return g.queue
}

// Controller returns the single controller assembled for g.
//
// A nil receiver returns nil. The controller consumes Queue and reads snapshots
// from Primary().Cache().
func (g *MultiSource) Controller() *objectcontroller.Controller {
	if g == nil {
		return nil
	}

	return g.controller
}
