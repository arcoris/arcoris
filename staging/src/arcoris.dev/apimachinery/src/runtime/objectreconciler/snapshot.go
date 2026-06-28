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

// Snapshot is the stable object read model passed to user reconciliation logic.
//
// Snapshot wraps the generic cache snapshot with object reconciliation
// vocabulary. Revision is the objectstore collection revision of View. The
// expected invariant for snapshots produced by ReconcileOnce is:
//
//	snap.Revision == snap.View.Revision()
type Snapshot struct {
	// Revision is the objectstore collection revision observed by this
	// reconciliation attempt.
	Revision objectstore.Revision

	// View is the detached latest object cache view for Revision.
	//
	// View is stable for the duration of user reconciliation and does not
	// observe later cache mutations.
	View objectcache.View
}

// snapshotFromSource converts a generic cache snapshot into reconciler input.
func snapshotFromSource(source snapshot.Snapshot[objectstore.Revision, objectcache.View]) Snapshot {
	return Snapshot{
		Revision: source.Revision,
		View:     source.Value,
	}
}
