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

// update commits report as the latest observation for target.
//
// The per-target snapshot.Store assigns Revision and Updated. The return value
// is false when target is not configured or report does not form a valid
// observation for target.
//
// The method validates and clones the observation before taking store.mu. That
// keeps the structural map lock focused on map ownership only. When the target
// already has a snapshot.Store, Replace is called while store.mu is held so the
// pointer cannot be removed or swapped between lookup and write. The per-target
// store still owns its own value lock and revision advance.
func (s *store) update(target health.Target, report health.Report) bool {
	obs, ok := newObservation(target, report)
	if !ok {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.configured[target]; !ok {
		return false
	}

	targetStore, ok := s.byTarget[target]
	if !ok {
		// NewStore commits the first observation immediately at revision 1 and
		// records Updated using the configured clock.
		s.byTarget[target] = snapshot.NewStore(
			obs,
			cloneObservation,
			snapshot.WithClock(s.clock),
		)
		return true
	}

	targetStore.Replace(obs)
	return true
}
