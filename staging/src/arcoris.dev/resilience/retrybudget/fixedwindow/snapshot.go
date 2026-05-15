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

package fixedwindow

import (
	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

// Snapshot returns the latest published retry-budget snapshot.
//
// Snapshot is lock-free. Window rotation is performed on write paths, not on
// read paths, so a quiet limiter may keep publishing the last window until the
// next RecordOriginal or TryAdmitRetry call observes time advancement.
func (l *Limiter) Snapshot() snapshot.Snapshot[retrybudget.Snapshot] {
	return l.published.Snapshot()
}

// Revision returns the latest published source-local revision.
func (l *Limiter) Revision() snapshot.Revision {
	return l.published.Revision()
}

// snapshotValueLocked builds the immutable domain snapshot for the current
// limiter state.
//
// The caller must hold l.mu.
func (l *Limiter) snapshotValueLocked() retrybudget.Snapshot {
	allowed := allowedRetries(l.original, l.cfg.ratio, l.cfg.minRetries)
	available := availableRetries(allowed, l.retries)

	return retrybudget.Snapshot{
		Kind: retrybudget.KindFixedWindow,
		Attempts: retrybudget.AttemptsSnapshot{
			Original: l.original,
			Retry:    l.retries,
		},
		Capacity: retrybudget.CapacitySnapshot{
			Allowed:   allowed,
			Available: available,
			Exhausted: available == 0,
		},
		Window: retrybudget.WindowSnapshot{
			StartedAt: l.windowStart,
			EndsAt:    l.windowEndLocked(),
			Duration:  l.cfg.window,
			Bounded:   true,
		},
		Policy: retrybudget.PolicySnapshot{
			Ratio:   l.cfg.ratio,
			Minimum: l.cfg.minRetries,
			Bounded: true,
		},
	}
}

// publishLocked publishes the current immutable domain snapshot.
//
// The caller must hold l.mu.
func (l *Limiter) publishLocked() snapshot.Snapshot[retrybudget.Snapshot] {
	return l.published.Publish(l.snapshotValueLocked())
}
