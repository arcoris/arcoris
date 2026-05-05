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

package healthprobe

import (
	"sync"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

// store is the concurrency-safe latest-snapshot cache owned by Runner.
//
// store is intentionally private. Runner exposes copy-safe read methods while
// retaining freedom to change cache internals later. store does not call
// Evaluator, compute staleness, start goroutines, or interpret health policy.
type store struct {
	mu sync.RWMutex

	targets     []health.Target
	byTarget    map[health.Target]Snapshot
	generations map[health.Target]uint64
}

// newStore returns an empty latest-snapshot store for targets.
func newStore(targets []health.Target) *store {
	return &store{
		targets:     copyTargets(targets),
		byTarget:    make(map[health.Target]Snapshot, len(targets)),
		generations: make(map[health.Target]uint64, len(targets)),
	}
}

// update stores report as the latest snapshot for target.
//
// update clones report before storing it so callers cannot mutate cached state
// through a slice returned by health.Report.Checks. Generation is incremented per
// target only after the report is accepted by the cache.
func (s *store) update(target health.Target, report health.Report, updated time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	generation := s.generations[target] + 1
	s.generations[target] = generation

	s.byTarget[target] = Snapshot{
		Target:     target,
		Report:     cloneReport(report),
		Updated:    updated,
		Generation: generation,
	}
}

// snapshot returns the latest snapshot for target.
//
// The returned snapshot is detached from internal cache slices. The boolean is
// false when target has not yet been observed or is not known to the store.
func (s *store) snapshot(target health.Target) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot, ok := s.byTarget[target]
	if !ok {
		return Snapshot{}, false
	}

	return cloneSnapshot(snapshot), true
}

// snapshots returns all observed snapshots in configured target order.
//
// Unobserved targets are omitted. Every returned snapshot is detached from
// internal cache slices.
func (s *store) snapshots() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots := make([]Snapshot, 0, len(s.byTarget))
	for _, target := range s.targets {
		snapshot, ok := s.byTarget[target]
		if !ok {
			continue
		}

		snapshots = append(snapshots, cloneSnapshot(snapshot))
	}

	return snapshots
}
