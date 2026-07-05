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
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// MultiSourceConfig describes one controller graph with multiple watched
// inputs feeding one shared queue.
//
// Primary is required and provides the controller snapshot source. Secondary
// inputs are optional; each secondary input has its own cache and reflector, but
// emits mapped work into the same queue as the primary input.
type MultiSourceConfig struct {
	// Primary is the watched input whose cache is used as the controller
	// snapshot source.
	Primary InputConfig

	// Secondary are additional watched inputs that feed the same shared queue.
	Secondary []InputConfig

	// Reconciler performs reconciliation attempts for requests from the shared
	// queue.
	Reconciler objectreconciler.Reconciler

	// Queue configures the shared bounded work queue.
	Queue objectworkqueue.Options

	// Controller configures the single controller consuming the shared queue.
	Controller objectcontroller.Options
}
