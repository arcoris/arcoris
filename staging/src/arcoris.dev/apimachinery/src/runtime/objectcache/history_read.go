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

// historicalReadState is the common validated state for historical reads.
//
// current is captured under Cache.mu so GetAt and PreviousLive compare against
// a stable cache boundary while scanning the selected object record.
type historicalReadState struct {
	current objectstore.Revision
	record  *objectRecord
}

// prepareHistoricalReadLocked validates common historical-read preconditions.
//
// Cache.mu must be held. The caller validates key collection membership before
// locking because collection is immutable after New.
func (c *Cache) prepareHistoricalReadLocked(key objectstore.Key, revision objectstore.Revision) (historicalReadState, error) {
	if !c.ready {
		return historicalReadState{}, ErrNotReady
	}
	if !c.historyEnabled {
		return historicalReadState{}, ErrHistoryDisabled
	}
	current := c.latest.revision
	if current.Before(revision) {
		return historicalReadState{}, ErrFutureRevision
	}

	return historicalReadState{current: current, record: c.records[key]}, nil
}

func knownAbsenceAt(key objectstore.Key, revision objectstore.Revision) GetResult {
	return GetResult{Key: key, Revision: revision}
}

// liveResultAt returns a detached live read result at the served cache
// observation revision.
func liveResultAt(key objectstore.Key, revision objectstore.Revision, state objectstore.State) GetResult {
	return GetResult{Key: key, State: state.Clone(), Found: true, Revision: revision}
}
