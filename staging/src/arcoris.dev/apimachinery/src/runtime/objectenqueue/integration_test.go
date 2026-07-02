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

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestIntegrationWithObjectWorkQueue(t *testing.T) {
	tests := []struct {
		name   string
		change objectstore.Change
	}{
		{name: "created", change: createdChange(t, 1)},
		{name: "updated", change: updatedChange(t, 2)},
		{name: "deleted", change: deletedChange(t, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 1})
			requireNoError(t, err)

			handler := mustHandler(t, queue, ChangedObject())
			requireNoError(t, handler.Handle(context.Background(), tt.change))

			item, err := queue.Get(context.Background())
			requireNoError(t, err)
			requireItem(t, item, objectworkqueue.Item{Key: tt.change.Key})
			requireNoError(t, queue.Done(item))
		})
	}
}
