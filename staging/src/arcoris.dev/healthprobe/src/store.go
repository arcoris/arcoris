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
	"sync"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
	"arcoris.dev/snapshot"
)

// store is the target-indexed latest-observation cache owned by Runner.
//
// store owns target configuration, configured target order, and the target to
// per-target store map. Each per-target snapshot.Store[observation] owns the
// latest observation value, revision, update timestamp, concurrency safety, and
// clone isolation for that target. store.mu protects only the map and target
// configuration structure.
//
// The split is intentional:
//   - store answers package-specific questions such as "is this target
//     configured?" and "what order should Snapshots return?";
//   - snapshot.Store answers generic state-holder questions such as "what is the
//     latest value?", "what revision is it?", "when was it committed?", and
//     "how do reads and writes stay detached?".
type store struct {
	// mu protects targets, configured, and byTarget. It must not be held while
	// calling Stamped or Replace on a per-target snapshot.Store; those stores own
	// their own locks and clone boundaries.
	mu sync.RWMutex

	// clock is passed to newly-created per-target snapshot stores so Updated is
	// assigned by the same clock Runner uses for stale calculations.
	clock clock.PassiveClock

	// targets preserves the configured target order used by Snapshots.
	targets []health.Target

	// configured is the membership set for fast rejection of unconfigured
	// updates and reads.
	configured map[health.Target]struct{}

	// byTarget contains one snapshot.Store per target after that target has its
	// first valid observation. Missing entries represent configured but
	// unobserved targets.
	byTarget map[health.Target]*snapshot.Store[observation]
}

// newStore returns an empty latest-observation store for targets.
//
// The constructor copies targets because Runner construction accepts caller-owned
// slices. The copy becomes the stable order used by snapshots. Per-target
// snapshot stores are intentionally created lazily: a configured target with no
// valid observation must still read as absent, not as a zero-valued committed
// observation.
func newStore(targets []health.Target, clk clock.PassiveClock) *store {
	copied := copyTargets(targets)
	configured := make(map[health.Target]struct{}, len(copied))
	for _, target := range copied {
		configured[target] = struct{}{}
	}

	return &store{
		clock:      clk,
		targets:    copied,
		configured: configured,
		byTarget:   make(map[health.Target]*snapshot.Store[observation], len(copied)),
	}
}
