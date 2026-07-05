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
	"arcoris.dev/apimachinery/runtime/objectenqueue"
)

// InputConfig describes one watched input collection in a multi-source graph.
//
// Each input owns its own cache and reflector, but emits work into the graph's
// shared queue. The input does not decide how those work items are reconciled;
// the single graph controller consumes the shared queue.
type InputConfig struct {
	// Source provides list-watch access for this watched input collection.
	Source storewatchapi.ListerWatcher

	// Collection identifies the watched input collection reflected into this
	// input's cache.
	Collection objectstore.ListRequest

	// Predicate selects which input objects and changes are considered for
	// mapping. The zero Predicate means all input objects and changes.
	Predicate objectquery.Predicate

	// Listed maps list and relist items to work items for the shared queue.
	Listed objectenqueue.ListItemMapper

	// Changed maps committed changes to work items for the shared queue.
	Changed objectenqueue.Mapper
}
