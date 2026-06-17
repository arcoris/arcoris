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
// Snapshot is stronger than raw objectstore.Store.List because it promises that
// the returned ListResult is tied to a boundary intended for a subsequent watch
// opened on the same ListerWatcher or Store boundary.
type Snapshotter interface {
	// Snapshot reads collection and returns a validated list-to-watch snapshot.
	Snapshot(context.Context, objectstore.ListRequest) (Snapshot, error)
}
