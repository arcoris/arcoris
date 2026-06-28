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

package objectreconciler

import (
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/snapshot"
)

// SnapshotSource provides detached object cache snapshots for reconciliation.
//
// SnapshotSource is intentionally specific to objectcache.View. Generic
// snapshot source interfaces live in arcoris.dev/snapshot; this package is the
// object reconciliation execution boundary.
type SnapshotSource interface {
	// ReadSnapshot returns the latest object cache view or an error such as
	// objectcache.ErrNotReady when the read model has not been initialized.
	//
	// A successful call must return a detached, stable view whose revision
	// matches the returned snapshot revision.
	ReadSnapshot() (snapshot.Snapshot[objectstore.Revision, objectcache.View], error)
}
