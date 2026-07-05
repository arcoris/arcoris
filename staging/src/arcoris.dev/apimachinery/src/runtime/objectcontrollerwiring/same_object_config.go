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
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// SameObjectConfig describes the dependencies and options needed to assemble
// one same-object controller graph.
//
// The config is intentionally declarative. Construction validates and wires the
// lower-level primitives, but does not start the reflector, run controller
// workers, or shut down the queue.
type SameObjectConfig struct {
	// Source provides list-watch access for Collection.
	//
	// NewSameObject passes Source directly to objectreflector.New. Source
	// remains the owner of collection continuity and watch-stream behavior.
	Source storewatchapi.ListerWatcher

	// Collection identifies the object collection reflected into the cache and
	// converted into same-object reconciliation work.
	Collection objectstore.ListRequest

	// Reconciler performs user reconciliation attempts for requests produced
	// from reflected object keys.
	Reconciler objectreconciler.Reconciler

	// Queue configures the bounded object work queue created for the graph.
	Queue objectworkqueue.Options

	// Controller configures the fixed controller workers created for the graph.
	Controller objectcontroller.Options

	// Predicate selects which reflected objects produce same-object work.
	// The zero Predicate means all objects.
	Predicate objectquery.Predicate

	// Indexes are optional prebuilt secondary indexes for Collection.
	//
	// NewSameObject installs indexes after the cache sink and before the enqueue
	// sink. Mappers can close over these same index instances when they need
	// read-side reverse lookups; SameObject itself still emits same-key work.
	Indexes []*objectindex.Index
}
