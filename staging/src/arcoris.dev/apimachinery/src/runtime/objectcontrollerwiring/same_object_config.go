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
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// SameObjectConfig describes one same-object controller component graph.
type SameObjectConfig struct {
	// Source provides list-watch access for the reflected collection.
	Source storewatchapi.ListerWatcher

	// Collection identifies the object collection watched by this controller.
	Collection objectstore.ListRequest

	// Reconciler performs user reconciliation attempts.
	Reconciler objectreconciler.Reconciler

	// Queue configures the bounded object work queue.
	Queue objectworkqueue.Options

	// Controller configures controller workers.
	Controller objectcontroller.Options

	// Predicate selects which reflected objects produce same-object work.
	// The zero Predicate means all objects.
	Predicate objectquery.Predicate
}
