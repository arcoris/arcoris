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

package objectcache

import "arcoris.dev/apimachinery/api/objectstore"

// ListResult is the result of querying a Snapshot.
//
// Revision is always the snapshot revision, not a revision recomputed from the
// matching items. Empty query results still carry the snapshot revision.
type ListResult struct {
	// Items are detached list items that matched the query.
	Items []objectstore.ListItem

	// Revision is the snapshot's observed store revision watermark.
	Revision objectstore.Revision
}
