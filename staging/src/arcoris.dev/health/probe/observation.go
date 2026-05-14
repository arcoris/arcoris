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

import (
	"arcoris.dev/health"
	"arcoris.dev/snapshot"
)

// observation is the private value stored in a per-target snapshot.Store.
//
// The type is deliberately smaller than public Snapshot. It contains only the
// domain payload that Runner wants to remember for one target: the target key and
// the health report observed for that key. It does not contain cache metadata.
//
// Revision and update time are assigned by snapshot.Store when the observation
// is committed. Keeping those fields out of observation prevents probe.store
// from accidentally reintroducing manual generation or timestamp ownership.
// Stale is also excluded because staleness is a read-boundary calculation: the
// same stored observation can be fresh at one read and stale at a later read.
type observation struct {
	// Target is copied from the store key so stamped observations can be adapted
	// back into public Snapshot values without trusting external map state.
	Target health.Target

	// Report is the root health domain report for Target. The report is cloned
	// before entering snapshot.Store and cloned again when leaving it.
	Report health.Report
}

// newObservation validates and clones the value that will cross into
// snapshot.Store ownership.
//
// The clone happens before validation so later caller mutations cannot affect the
// value that passed validation. The boolean keeps store.update small: false
// means the report is not a valid observed value for target and must not advance
// the per-target snapshot revision.
func newObservation(target health.Target, report health.Report) (observation, bool) {
	obs := observation{
		Target: target,
		Report: cloneReport(report),
	}
	if !obs.isObserved() {
		return observation{}, false
	}

	return obs, true
}

// isObserved reports whether o is a complete cache payload.
//
// This predicate intentionally checks only payload invariants. It does not check
// Revision, Updated, or Stale because those fields are not part of observation.
// Public Snapshot performs the full read-model validation after snapshot.Store
// has added its stamped metadata.
func (o observation) isObserved() bool {
	return o.Target.IsConcrete() &&
		o.Report.Target == o.Target &&
		o.Report.IsValid()
}

// snapshotFromStamped adapts the snapshot package read model to probe.Snapshot.
//
// snapshot.Store returns a stamped observation: payload plus source-local
// Revision and commit Updated time. probe.Snapshot is the public domain read
// model, so this function copies the payload fields across and preserves the
// stamped metadata. Stale is intentionally left false; Runner sets it at the
// read boundary after comparing Updated with the configured stale-after window.
func snapshotFromStamped(stamped snapshot.Stamped[observation]) Snapshot {
	obs := stamped.Value

	return Snapshot{
		Target:   obs.Target,
		Report:   cloneReport(obs.Report),
		Revision: stamped.Revision,
		Updated:  stamped.Updated,
	}
}
