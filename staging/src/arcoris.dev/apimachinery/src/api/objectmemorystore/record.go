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

package objectmemorystore

import "arcoris.dev/apimachinery/api/objectstore"

// record is an immutable published object record.
//
// Live records hold committed object state. Tombstone records remember the
// delete commit revision while keeping the previous live state available to the
// delete operation that won the compare-and-swap.
type record struct {
	// state is detached committed live state. It must never be mutated.
	state objectstore.State

	// deleteRevision is non-zero only for tombstone records.
	deleteRevision objectstore.Revision

	// deleted marks this record as a tombstone.
	deleted bool
}
