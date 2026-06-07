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
	"time"

	"arcoris.dev/health"
	"arcoris.dev/snapshot"
)

// Snapshot is the public read model for the latest probed health state of one
// target.
//
// A Snapshot is produced by Runner after a scheduled health evaluation and is
// returned to callers through Runner read methods. It is intentionally a plain,
// copyable value: callers may store it, compare it, render it, or adapt it
// without owning Runner internals.
//
// Snapshot does not define transport behavior. HTTP status mapping, gRPC serving
// state, metrics, logging, restart policy, admission policy, routing, and
// scheduler decisions belong outside package probe.
//
// Snapshot is built from the stamped observation held by the internal per-target
// snapshot.Store. Stored reports must be structurally valid and aggregate
// consistent; probe rejects stale aggregate reports instead of repairing them at
// the cache boundary. The cache surface stays report-oriented and does not
// expose evaluator execution errors. Transport adapters MUST still avoid
// exposing health.Result.Cause unless they explicitly own a safe diagnostic
// surface.
type Snapshot struct {
	// Target is the health target represented by this snapshot.
	//
	// Target duplicates Report.Target intentionally. The field lets consumers
	// inspect the cache key without having to trust or dereference the embedded
	// report first. Runner and store code must keep Target and Report.Target
	// consistent for observed snapshots.
	Target health.Target

	// Report is the latest target-level health report stored for Target.
	//
	// Report.Checks must be defensively copied whenever a Snapshot crosses the
	// store boundary. health.Report is a plain value, but its Checks field is a
	// slice and therefore has shared backing-array semantics unless copied.
	Report health.Report

	// Revision is the per-target snapshot store revision.
	//
	// Revision advances each time Runner commits an observation for Target. It is
	// local to one Runner and one target. It is not a global ordering,
	// distributed resource version, persistence version, lease epoch, or fencing
	// token.
	Revision snapshot.Revision

	// Updated is the time at which the per-target snapshot store committed
	// Report.
	//
	// Updated is distinct from Report.Observed. Report.Observed belongs to the
	// health evaluation itself; Updated belongs to the snapshot store commit
	// boundary. In normal operation they are close, but they intentionally
	// describe different events.
	Updated time.Time

	// Stale reports whether the snapshot was older than the configured staleness
	// window when it was read.
	//
	// Stale is computed at the read boundary by Runner. It is not a stored fact: a
	// fresh snapshot can become stale without a new write. Callers should treat
	// Stale as read-time cache metadata, not as part of health.Report.
	Stale bool
}
