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

package objectcontroller

import (
	"context"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

var _ Queue = (*recordingQueue)(nil)

func TestQueueContractUsesOnlyGetAndDone(t *testing.T) {
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1)}}

	item, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireNoError(t, queue.Done(item))

	if queue.getCallCount() != 1 || queue.doneCount() != 1 {
		t.Fatalf("get calls=%d done calls=%d; want 1,1", queue.getCallCount(), queue.doneCount())
	}
}

type recordingQueue struct {
	mu        sync.Mutex
	items     []objectworkqueue.Item
	getErrs   []error
	doneErr   error
	getCalls  int
	doneItems []objectworkqueue.Item
}

func (q *recordingQueue) Get(context.Context) (objectworkqueue.Item, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.getCalls++
	if len(q.getErrs) > 0 {
		err := q.getErrs[0]
		q.getErrs = q.getErrs[1:]
		if err != nil {
			return objectworkqueue.Item{}, err
		}
	}
	if len(q.items) == 0 {
		return objectworkqueue.Item{}, objectworkqueue.ErrShutDown
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *recordingQueue) Done(item objectworkqueue.Item) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.doneItems = append(q.doneItems, item)
	return q.doneErr
}

func (q *recordingQueue) doneCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.doneItems)
}

func (q *recordingQueue) completedItems() []objectworkqueue.Item {
	q.mu.Lock()
	defer q.mu.Unlock()

	return append([]objectworkqueue.Item(nil), q.doneItems...)
}

func (q *recordingQueue) getCallCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.getCalls
}
