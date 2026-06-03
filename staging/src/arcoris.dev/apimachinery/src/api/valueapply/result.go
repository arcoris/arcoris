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

package valueapply

import (
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuecompare"
)

// Result contains the merged value, updated ownership, and apply metadata.
//
// Early validation errors may return a zero Result. Conflict errors return the
// field sets and comparison computed before conflict handling. Merge errors
// return all pre-merge metadata, but not post-merge ownership.
type Result struct {
	// Value is the merged value.
	Value value.Value

	// Ownership is the updated field ownership state.
	Ownership fieldownership.State

	// AppliedFields are semantic fields mentioned by Applied.
	AppliedFields fieldpath.Set

	// DroppedFields are fields previously owned by Owner but omitted by Applied.
	DroppedFields fieldpath.Set

	// DeletedFields are dropped fields actually removed from Value because no
	// other owner overlapped them.
	DeletedFields fieldpath.Set

	// ChangedAppliedFields are applied fields whose value would change Live.
	ChangedAppliedFields fieldpath.Set

	// MergeFields are fields passed to valuemerge.
	MergeFields fieldpath.Set

	// Changes is the live/applied semantic comparison.
	Changes valuecompare.Result

	// Conflicts are ownership conflicts found before force handling.
	Conflicts fieldownership.ConflictSet
}
