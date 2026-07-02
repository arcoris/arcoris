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

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// handlerEmitter adapts Handler queue state into the emit callback passed to a Mapper.
type handlerEmitter struct {
	// ctx is forwarded unchanged to Enqueuer.Add.
	ctx context.Context

	// queue receives every mapped objectworkqueue.Item.
	queue Enqueuer

	// err stores the first queue error so later emit calls return the same error.
	err error
}

// emit forwards item to the queue unless a previous emit already failed.
func (e *handlerEmitter) emit(item objectworkqueue.Item) error {
	if e.err != nil {
		return e.err
	}

	e.err = e.queue.Add(e.ctx, item)

	return e.err
}

// result gives queue errors precedence over mapper return errors.
func (e *handlerEmitter) result(mapErr error) error {
	if e.err != nil {
		return e.err
	}

	return mapErr
}
