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

package objectenqueue

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Handler maps committed object changes and enqueues the resulting work items.
type Handler struct {
	// queue is the producer-side enqueue target used by emitted work items.
	queue Enqueuer

	// mapper owns the change-to-work mapping policy for this handler.
	mapper Mapper
}

// New constructs a Handler from an enqueue target and a mapper.
//
// New validates only local wiring. It does not validate future changes, start
// workers, consume watches, or touch queue state.
func New(queue Enqueuer, mapper Mapper) (*Handler, error) {
	if isNilInterface(queue) {
		return nil, ErrNilQueue
	}
	if isNilInterface(mapper) {
		return nil, ErrNilMapper
	}

	return &Handler{queue: queue, mapper: mapper}, nil
}

// Handle maps change and forwards each emitted item to the handler queue.
//
// Handle keeps no mutable per-call state on Handler. If the configured queue
// and mapper are safe for concurrent use, concurrent Handle calls are safe too.
// The ctx value is passed through to Enqueuer.Add unchanged.
func (h *Handler) Handle(ctx context.Context, change objectstore.Change) error {
	if h == nil || isNilInterface(h.queue) || isNilInterface(h.mapper) {
		return ErrInvalidHandler
	}

	emitter := handlerEmitter{ctx: ctx, queue: h.queue}
	mapErr := h.mapper.Map(change, emitter.emit)

	return emitter.result(mapErr)
}
