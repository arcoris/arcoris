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

// IsZero reports whether s is the zero snapshot value.
//
// The zero value represents "no cached observation." Runner read methods should
// normally return ok=false instead of returning a zero Snapshot, but IsZero is
// useful in tests and defensive integration code.
//
// A zero Snapshot must also have snapshot.ZeroRevision. A committed
// snapshot.Store observation can never have ZeroRevision, so this keeps absence
// distinguishable from a real cached observation.
func (s Snapshot) IsZero() bool {
	return s.Target == health.TargetUnknown &&
		reportIsZero(s.Report) &&
		s.Revision == snapshot.ZeroRevision &&
		s.Updated.IsZero() &&
		!s.Stale
}

// IsObserved reports whether s contains a stored probe observation.
//
// IsObserved is intentionally stricter than checking Updated and Revision
// alone. A snapshot is observed only when the complete Snapshot invariants hold:
// a concrete Target, a valid Report for the same Target, a non-zero cache update
// timestamp, and a non-zero revision. The embedded Report may still represent
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
// Report is structurally valid, Report is aggregate-consistent, Report.Target
// matches Target, Revision is non-zero, and Updated is non-zero. Stale may be
// true only on an otherwise observed Snapshot because stale is read-time cache
// metadata.
//
// IsValid intentionally does not interpret health success or failure. A valid
// Snapshot may contain a healthy, degraded, unhealthy, or unknown health.Report.
// This method checks only whether the cache read model is structurally safe to
// hand to callers.
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
	if !s.Report.IsConsistent() {
		return false
	}
	if s.Revision == snapshot.ZeroRevision {
		return false
	}
	if s.Updated.IsZero() {
		return false
	}

	return true
}

// reportIsZero reports whether report is the zero health.Report value.
//
// health.Report is not directly comparable because it contains the Checks slice.
// Keep this helper local to package probe instead of changing package health only
// to support Snapshot.IsZero.
//
// A nil Checks slice and an empty Checks slice are both treated as zero here
// because a zero report means "no report payload", not a committed report with an
// intentionally empty check list. Committed reports are validated through
// health.Report.IsValid instead.
func reportIsZero(report health.Report) bool {
	return report.Target == health.TargetUnknown &&
		report.Status == health.StatusUnknown &&
		report.Observed.IsZero() &&
		report.Duration == 0 &&
		len(report.Checks) == 0
}
