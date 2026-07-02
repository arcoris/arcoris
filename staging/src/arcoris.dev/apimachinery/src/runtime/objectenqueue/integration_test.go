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
	"testing"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestIntegrationWithObjectWorkQueue(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 1})
	requireNoError(t, err)

	change := createdChange(t, 1)
	handler := mustHandler(t, queue, ChangedObject())

	requireNoError(t, handler.Handle(context.Background(), change))

	item, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, item, objectworkqueue.Item{Key: change.Key})
	requireNoError(t, queue.Done(item))
}
