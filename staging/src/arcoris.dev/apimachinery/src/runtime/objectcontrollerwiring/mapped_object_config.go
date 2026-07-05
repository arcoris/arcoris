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
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// MappedObjectConfig describes the dependencies and mapping policy needed to
// assemble one mapped-object controller graph.
//
// The graph watches Collection, stores that source collection in its cache, and
// lets Listed and Changed emit reconciliation work for arbitrary mapped keys.
// The reconciler receives mapped requests together with a snapshot of the
// watched source collection.
type MappedObjectConfig struct {
	// Source provides list-watch access for the watched source collection.
	Source storewatchapi.ListerWatcher

	// Collection identifies the watched source collection.
	Collection objectstore.ListRequest

	// Reconciler performs reconciliation attempts for mapped target requests.
	Reconciler objectreconciler.Reconciler

	// Queue configures the bounded object work queue created for mapped work.
	Queue objectworkqueue.Options

	// Controller configures the fixed controller workers created for the graph.
	Controller objectcontroller.Options

	// Predicate selects which source objects and changes are considered for
	// mapping. The zero Predicate means all source objects and changes.
	Predicate objectquery.Predicate

	// Listed maps source list items from Replace/relist boundaries to target
	// work items.
	Listed objectenqueue.ListItemMapper

	// Changed maps committed source changes to target work items.
	Changed objectenqueue.Mapper
}
