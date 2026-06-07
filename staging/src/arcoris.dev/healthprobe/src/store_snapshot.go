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

package probe

import (
	"arcoris.dev/health"
	"arcoris.dev/snapshot"
)

// snapshot returns the latest observed snapshot for target.
//
// The returned snapshot is detached from internal store slices. The boolean is
// false when target has not yet been observed or is not known to the store.
//
// store.mu is released before calling Stamped. The target store pointer is
// stable after lookup, and snapshot.Store owns concurrent read/write safety for
// the actual observation. Releasing store.mu avoids nesting the structural map
// lock around clone work and per-target store locking.
func (s *store) snapshot(target health.Target) (Snapshot, bool) {
	s.mu.RLock()
	if _, ok := s.configured[target]; !ok {
		s.mu.RUnlock()
		return Snapshot{}, false
	}
	targetStore := s.byTarget[target]
	s.mu.RUnlock()

	if targetStore == nil {
		return Snapshot{}, false
	}

	snap := snapshotFromStamped(targetStore.Stamped())
	if !snap.IsObserved() {
		return Snapshot{}, false
	}

	return snap, true
}

// snapshots returns all observed snapshots in configured target order.
//
// Unobserved targets are omitted. Every returned snapshot is detached from
// internal store slices.
//
// The method copies both configured order and target-store pointers while holding
// store.mu, then performs stamped reads after unlocking. That preserves a
// consistent order snapshot without blocking updates on report cloning.
func (s *store) snapshots() []Snapshot {
	s.mu.RLock()
	targets := copyTargets(s.targets)
	stores := make(map[health.Target]*snapshot.Store[observation], len(s.byTarget))
	for target, targetStore := range s.byTarget {
		stores[target] = targetStore
	}
	s.mu.RUnlock()

	snapshots := make([]Snapshot, 0, len(stores))
	for _, target := range targets {
		targetStore := stores[target]
		if targetStore == nil {
			continue
		}

		snap := snapshotFromStamped(targetStore.Stamped())
		if snap.IsObserved() {
			snapshots = append(snapshots, snap)
		}
	}

	return snapshots
}
