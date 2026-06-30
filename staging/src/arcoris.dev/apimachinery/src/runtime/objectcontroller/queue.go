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

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// Queue is the read/complete subset of an object work queue consumed by a
// Controller.
//
// Controller workers call Get to receive work and Done exactly once after each
// successful Get. The producer side of a queue is intentionally outside this
// package.
type Queue interface {
	// Get returns the next queued item or an error.
	//
	// objectworkqueue.ErrShutDown means the queue is shut down and drained, so
	// the worker exits cleanly. Other errors are fatal to Run.
	Get(context.Context) (objectworkqueue.Item, error)

	// Done completes processing for item.
	//
	// Done errors are fatal to Run, but the controller still guarantees a
	// single Done attempt for every item returned by Get.
	Done(objectworkqueue.Item) error
}
