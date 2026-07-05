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
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// MappedObjectConfig describes the dependencies and mapping policy for one
// mapped-object controller graph.
//
// The graph watches Collection, stores that source collection in its cache, and
// lets Listed and Changed emit reconciliation work for arbitrary mapped keys.
// The reconciler receives mapped requests together with a snapshot of the
// watched source collection. Construction validates by delegating to the
// lower-level primitives; this package does not reinterpret their option or
// dependency errors.
type MappedObjectConfig struct {
	// Source provides list-watch access for the watched source collection.
	//
	// Source is passed directly to objectreflector.New. It remains responsible
	// for list/watch continuity, retained history, and watch stream behavior.
	Source storewatchapi.ListerWatcher

	// Collection identifies the watched source collection reflected into Cache.
	//
	// MappedObject does not infer or validate any target collection from mapped
	// output keys. The queue accepts objectworkqueue.Item values emitted by the
	// configured mappers.
	Collection objectstore.ListRequest

	// Reconciler performs reconciliation attempts for mapped target requests.
	//
	// Reconciler receives objectreconciler.Request values produced by the queue
	// and snapshots of Collection. Any additional target-state dependencies are
	// explicit reconciler dependencies, not part of this wiring helper.
	Reconciler objectreconciler.Reconciler

	// Queue configures the bounded object work queue created for mapped work.
	//
	// The queue bounds distinct pending mapped requests, regardless of how many
	// source objects emitted those requests.
	Queue objectworkqueue.Options

	// Controller configures the fixed controller workers created for the graph.
	Controller objectcontroller.Options

	// Predicate selects which source objects and changes are considered for
	// mapping. The zero Predicate means all source objects and changes.
	Predicate objectquery.Predicate

	// Listed maps source list items from Replace/relist boundaries to target
	// work items.
	//
	// Listed is used only by objectenqueue.ReflectorSink.Replace. It may emit
	// zero, one, or many work items for each listed source object.
	Listed objectenqueue.ListItemMapper

	// Changed maps committed source changes to target work items.
	//
	// Changed is used only by objectenqueue.ReflectorSink.ApplyChange. Leaving a
	// query may still need to emit cleanup work; that decision belongs to the
	// mapper and objectenqueue projection semantics.
	Changed objectenqueue.Mapper

	// Indexes are optional prebuilt secondary indexes for Collection.
	//
	// NewMappedObject installs indexes between cache and enqueue. Caller-provided
	// mappers may close over the same index instances and call Lookup while
	// mapping; this package never interprets index names or values.
	Indexes []*objectindex.Index
}
