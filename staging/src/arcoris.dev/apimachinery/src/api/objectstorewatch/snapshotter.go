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

package objectstorewatch

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Snapshotter produces boundary-safe snapshots for structural collections.
//
// Snapshot is not merely objectstore.Store.List under another name. It returns
// a validated list result tied to a revision boundary that is intended for
// subsequent watch continuation on the same ListerWatcher or Store boundary.
// Snapshotter alone does not prove that future history is retained; that proof
// happens when the same objectwatch.Source accepts and serves the watch request.
type Snapshotter interface {
	// Snapshot reads collection and returns a validated list-to-watch snapshot.
	Snapshot(context.Context, objectstore.ListRequest) (Snapshot, error)
}
