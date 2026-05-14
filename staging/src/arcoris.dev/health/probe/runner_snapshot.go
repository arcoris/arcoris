/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package probe

import "arcoris.dev/health"

// Snapshot returns the latest observed snapshot for target.
//
// The returned snapshot is detached from Runner internals. The boolean is false
// when the Runner is nil, target is not concrete, target is not configured, or no
// valid observation has been stored for target yet. Stale is computed at the
// read boundary using Runner's configured clock.
//
// store.snapshot returns an unstaled public Snapshot built from
// snapshot.Stamped[observation]. Runner then adds Stale because staleness depends
// on the read time and stale-after window, not only on the committed value.
func (r *Runner) Snapshot(target health.Target) (Snapshot, bool) {
	if r == nil || !target.IsConcrete() || !containsTarget(r.targets, target) {
		return Snapshot{}, false
	}

	snapshot, ok := r.store.snapshot(target)
	if !ok {
		return Snapshot{}, false
	}

	snapshot = r.withReadStale(snapshot)
	if !snapshot.IsObserved() {
		return Snapshot{}, false
	}

	return snapshot, true
}

// Snapshots returns all observed snapshots in configured target order.
//
// A nil Runner returns nil. Unobserved or invalid snapshots are omitted. Each
// returned snapshot is detached from Runner internals. Stale is computed at the
// read boundary for each snapshot.
//
// The store preserves configured target order and returns only observed targets.
// Runner performs a final IsObserved check after stale calculation so invalid
// data cannot cross the public read boundary even if future store internals
// change.
func (r *Runner) Snapshots() []Snapshot {
	if r == nil {
		return nil
	}

	stored := r.store.snapshots()
	snapshots := make([]Snapshot, 0, len(stored))
	for _, snapshot := range stored {
		snapshot = r.withReadStale(snapshot)
		if snapshot.IsObserved() {
			snapshots = append(snapshots, snapshot)
		}
	}

	return snapshots
}

func (r *Runner) withReadStale(snapshot Snapshot) Snapshot {
	// Stale is not stored in observation or snapshot.Store. It is recomputed on
	// every read so a snapshot can age naturally without forcing a write.
	if !snapshot.IsObserved() {
		snapshot.Stale = false
		return snapshot
	}

	snapshot.Stale = isStale(r.clock.Since(snapshot.Updated), r.staleAfter)
	return snapshot
}
