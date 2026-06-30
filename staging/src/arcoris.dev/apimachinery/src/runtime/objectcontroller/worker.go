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
	"errors"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// runWorker processes queue items until the queue shuts down or a fatal error
// occurs.
func (c *Controller) runWorker(ctx context.Context) error {
	for {
		item, err := c.queue.Get(ctx)
		if errors.Is(err, objectworkqueue.ErrShutDown) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := c.processItem(ctx, item); err != nil {
			return err
		}
	}
}
