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

	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// processItem reconciles one item and attempts Done exactly once.
func (c *Controller) processItem(ctx context.Context, item objectworkqueue.Item) (err error) {
	defer func() {
		doneErr := c.queue.Done(item)
		if err == nil {
			err = doneErr
			return
		}
		if doneErr != nil {
			err = errors.Join(err, doneErr)
		}
	}()

	return objectreconciler.ReconcileOnce(ctx, c.source, c.reconciler).Err
}
