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
	"time"

	"arcoris.dev/health"
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
// Snapshot also does not expose evaluator execution errors. Runner converts
// evaluator observations into health.Report values and keeps the public cache
// surface report-oriented. Transport adapters MUST still avoid exposing
// health.Result.Cause unless they explicitly own a safe diagnostic surface.
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

	// Updated is the time at which Runner stored Report in its latest-snapshot
	// cache.
	//
	// Updated is distinct from Report.Observed. Report.Observed belongs to the
	// health evaluation itself; Updated belongs to the probe cache write boundary.
	// In normal operation they are close, but they intentionally describe
	// different events.
	Updated time.Time

	// Generation is the per-target cache generation.
	//
	// Runner increments Generation each time it stores a new snapshot for the same
	// target. Generation is local to one Runner instance and is not a distributed
	// resource version, fencing token, lease epoch, or persistence contract.
	Generation uint64

	// Stale reports whether the snapshot was older than the configured staleness
	// window when it was read.
	//
	// Stale is computed at the read boundary by Runner. It is not a stored fact: a
	// fresh snapshot can become stale without a new write. Callers should treat
	// Stale as read-time cache metadata, not as part of health.Report.
	Stale bool
}

// IsZero reports whether s is the zero snapshot value.
//
// The zero value represents "no cached observation." Runner read methods should
// normally return ok=false instead of returning a zero Snapshot, but IsZero is
// useful in tests and defensive integration code.
func (s Snapshot) IsZero() bool {
	return s.Target == health.TargetUnknown &&
		reportIsZero(s.Report) &&
		s.Updated.IsZero() &&
		s.Generation == 0 &&
		!s.Stale
}

// IsObserved reports whether s contains a stored probe observation.
//
// IsObserved is intentionally stricter than checking Updated and Generation
// alone. A snapshot is observed only when the complete Snapshot invariants hold:
// a concrete Target, a valid Report for the same Target, a non-zero cache update
// timestamp, and a positive generation. The embedded Report may still represent
// an unknown health state.
func (s Snapshot) IsObserved() bool {
	return !s.IsZero() && s.IsValid()
}

// IsFresh reports whether s contains an observed snapshot that was not stale at
// the read boundary.
//
// Freshness is cache freshness, not health success. A fresh snapshot may still
// contain an unhealthy, degraded, or unknown health.Report.
func (s Snapshot) IsFresh() bool {
	return s.IsObserved() && !s.Stale
}

// IsValid reports whether s satisfies the Snapshot structural invariants.
//
// The zero Snapshot is valid and means that no cached observation exists. Any
// non-zero Snapshot must be a complete observed cache value: Target is concrete,
// Report is valid, Report.Target matches Target, Updated is non-zero, and
// Generation is positive. Stale may be true only on an otherwise observed
// Snapshot because stale is read-time cache metadata.
func (s Snapshot) IsValid() bool {
	if s.IsZero() {
		return true
	}
	if !s.Target.IsConcrete() {
		return false
	}
	if s.Report.Target != s.Target {
		return false
	}
	if !s.Report.IsValid() {
		return false
	}
	if s.Updated.IsZero() {
		return false
	}
	if s.Generation == 0 {
		return false
	}

	return true
}

// reportIsZero reports whether report is the zero health.Report value.
//
// health.Report is not directly comparable because it contains the Checks slice.
// Keep this helper local to package probe instead of changing package health only
// to support Snapshot.IsZero.
func reportIsZero(report health.Report) bool {
	return report.Target == health.TargetUnknown &&
		report.Status == health.StatusUnknown &&
		report.Observed.IsZero() &&
		report.Duration == 0 &&
		len(report.Checks) == 0
}
