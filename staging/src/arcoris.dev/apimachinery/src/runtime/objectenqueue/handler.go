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
	"reflect"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// Handler maps committed object changes and enqueues the resulting work items.
type Handler struct {
	queue  Enqueuer
	mapper Mapper
}

// New constructs a Handler from an enqueue target and a mapper.
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
func (h *Handler) Handle(ctx context.Context, change objectstore.Change) error {
	if h == nil || isNilInterface(h.queue) || isNilInterface(h.mapper) {
		return ErrInvalidHandler
	}

	var emitErr error
	emit := func(item objectworkqueue.Item) error {
		if emitErr != nil {
			return emitErr
		}
		emitErr = h.queue.Add(ctx, item)
		return emitErr
	}

	mapErr := h.mapper.Map(change, emit)
	if emitErr != nil {
		return emitErr
	}

	return mapErr
}

func isNilInterface(v any) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
