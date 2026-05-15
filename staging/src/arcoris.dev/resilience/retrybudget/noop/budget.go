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

package noop

import (
	"math"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

// Budget is an unlimited retry budget.
//
// Budget implements retrybudget.Budget but never denies retry admission.
// RecordOriginal is intentionally a no-op, TryAdmitRetry always returns an
// allowed decision, and Snapshot always returns the same immutable state.
//
// Budget is stateless, zero-value usable, and safe for concurrent use. It may be
// copied freely because it owns no mutable synchronization state.
type Budget struct{}

// New returns an unlimited retry budget.
func New() Budget {
	return Budget{}
}

// RecordOriginal records one original, non-retry attempt.
//
// For Budget this method is intentionally a no-op. The implementation does not
// track traffic because it never denies retries.
func (Budget) RecordOriginal() {}

// TryAdmitRetry admits one retry attempt.
//
// Budget always returns an allowed decision with a stable unlimited snapshot. The
// method does not record retry traffic.
func (Budget) TryAdmitRetry() retrybudget.Decision {
	return retrybudget.Decision{
		Allowed:  true,
		Reason:   retrybudget.ReasonAllowed,
		Snapshot: staticSnapshot(),
	}
}

// Snapshot returns the current unlimited retry-budget snapshot.
func (Budget) Snapshot() snapshot.Snapshot[retrybudget.Snapshot] {
	return staticSnapshot()
}

// Revision returns the stable revision used by the noop snapshot source.
func (Budget) Revision() snapshot.Revision {
	return staticRevision()
}

// staticSnapshot returns the immutable retry-budget snapshot exposed by Budget.
func staticSnapshot() snapshot.Snapshot[retrybudget.Snapshot] {
	return snapshot.Snapshot[retrybudget.Snapshot]{
		Revision: staticRevision(),
		Value:    staticSnapshotValue(),
	}
}

// staticRevision returns the source-local revision for Budget.
//
// The noop implementation has one always-present immutable state, so revision 1
// is used as its committed source-local revision.
func staticRevision() snapshot.Revision {
	return snapshot.ZeroRevision.Next()
}

// staticSnapshotValue returns the immutable domain value for Budget.
func staticSnapshotValue() retrybudget.Snapshot {
	return retrybudget.Snapshot{
		Attempts: retrybudget.AttemptsSnapshot{},
		Kind:     retrybudget.KindNoop,
		Capacity: retrybudget.CapacitySnapshot{
			Allowed:   math.MaxUint64,
			Available: math.MaxUint64,
			Exhausted: false,
		},
		Window: retrybudget.WindowSnapshot{
			Bounded: false,
		},
		Policy: retrybudget.PolicySnapshot{
			Bounded: false,
		},
	}
}

var _ retrybudget.Budget = Budget{}
var _ snapshot.Source[retrybudget.Snapshot] = Budget{}
var _ snapshot.RevisionSource = Budget{}
