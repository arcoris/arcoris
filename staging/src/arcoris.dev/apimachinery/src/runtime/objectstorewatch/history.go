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
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// changeHistory retains a bounded revision-ordered suffix of committed changes.
//
// Store.mu protects every field. The history is intentionally a simple slice:
// correctness and deterministic replay order matter more than indexing in this
// first passive wrapper.
type changeHistory struct {
	// changes is ordered by committed revision because Store serializes writes.
	changes []objectstore.Change
	// compactedThrough is the greatest revision that may have been discarded or
	// whose earlier continuity cannot be proven. A watch start before this
	// boundary cannot be served safely.
	compactedThrough objectstore.Revision
	// tainted records an unbounded history hole whose committed revision could
	// not be trusted. Historical starts fail explicitly while tainted.
	tainted bool
	// taintErr is the lower cause that made historical continuity unknowable.
	taintErr error
	// max is the retained change limit configured at construction.
	max int
}

// newChangeHistory constructs bounded retained history.
func newChangeHistory(max int) changeHistory {
	return changeHistory{max: max}
}

// append stores a detached change and compacts oldest retained items if needed.
//
// The caller must hold Store.mu. The input change is cloned so retained history
// cannot be mutated through caller-owned values.
func (h *changeHistory) append(change objectstore.Change) {
	h.changes = append(h.changes, change.Clone())
	for len(h.changes) > h.max {
		dropped := h.changes[0]
		copy(h.changes, h.changes[1:])
		last := len(h.changes) - 1
		// Clear the vacated slot so compacted object graphs are not retained by
		// the backing array after the slice is shortened.
		h.changes[last] = objectstore.Change{}
		h.changes = h.changes[:last]
		if h.compactedThrough.Before(dropped.Revision) {
			h.compactedThrough = dropped.Revision
		}
	}
}

// invalidateThrough marks history through revision as unsafe for replay.
//
// The caller must hold Store.mu. Retained changes at or before revision are
// discarded because a committed mutation at that boundary could not be
// represented as a valid watch event.
func (h *changeHistory) invalidateThrough(revision objectstore.Revision) {
	if !revision.IsValid() {
		return
	}
	if h.compactedThrough.Before(revision) {
		h.compactedThrough = revision
	}

	drop := 0
	for drop < len(h.changes) && !revision.Before(h.changes[drop].Revision) {
		drop++
	}
	if drop == 0 {
		return
	}

	copy(h.changes, h.changes[drop:])
	kept := len(h.changes) - drop
	for i := kept; i < len(h.changes); i++ {
		h.changes[i] = objectstore.Change{}
	}
	h.changes = h.changes[:kept]
}

// taint marks historical replay unsafe when the failed commit revision is
// unknown or invalid.
//
// The caller must hold Store.mu. Tainting does not prevent StartAtCurrent from
// observing a fresh backend boundary, but every historical start fails
// explicitly because no revision can prove it does not cross the unknown hole.
func (h *changeHistory) taint(cause error) {
	h.tainted = true
	if cause == nil {
		cause = errors.New("object store watch history tainted")
	}
	h.taintErr = cause
	h.changes = clearChanges(h.changes)
}

// replay returns cloned retained changes newer than revision for collection.
//
// The caller must hold Store.mu. The result is safe to enqueue into a stream
// without exposing history storage. A request older than compactedThrough fails
// explicitly because this wrapper cannot prove the missing prefix.
func (h changeHistory) replay(
	collection objectstore.ListRequest,
	revision objectstore.Revision,
) ([]objectstore.Change, error) {
	if h.tainted {
		return nil, objectwatch.HistoryUnavailable(h.taintErr)
	}
	if revision.Before(h.compactedThrough) {
		return nil, objectwatch.HistoryUnavailable(
			fmt.Errorf("requested revision %s is before compacted revision %s", revision, h.compactedThrough),
		)
	}

	var replay []objectstore.Change
	for _, change := range h.changes {
		if !revision.Before(change.Revision) {
			continue
		}
		if !objectstore.ChangeMatchesListRequest(change, collection) {
			continue
		}
		replay = append(replay, change.Clone())
	}

	return replay, nil
}

// clearChanges releases retained object graphs and returns an empty history slice.
func clearChanges(changes []objectstore.Change) []objectstore.Change {
	for i := range changes {
		changes[i] = objectstore.Change{}
	}

	return changes[:0]
}
