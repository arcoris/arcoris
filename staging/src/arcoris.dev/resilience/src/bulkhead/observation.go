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

package bulkhead

import "arcoris.dev/snapshot"

// Observation is direct metadata for one bulkhead acquisition attempt.
//
// Snapshot describes the observed local capacity state associated with the
// attempt. Refusal is none for successful acquisition and non-none for ordinary
// capacity denial. Observation intentionally carries no queues, waiters,
// fairness state, metrics labels, transport data, or historical counters.
//
// Observation is a value object. It is safe to store by value, pass through
// adapters, or compare in tests. It is not an event stream or global ordering
// token; under concurrent acquire, release, or limit changes, the snapshot is a
// diagnostic capacity observation rather than a serialization barrier for
// unrelated operations.
type Observation struct {
	// Snapshot is the revisioned capacity observation for the attempt.
	//
	// The snapshot comes from capacity.Ledger's explicit observed accounting
	// path. It exposes only scalar capacity state: Limit, Reserved, Available,
	// and Debt.
	Snapshot snapshot.Snapshot[Snapshot]

	// Refusal classifies why the attempt did not reserve capacity.
	//
	// Successful acquisition uses RefusalNone. Denied acquisition uses a
	// non-none capacity refusal such as RefusalInsufficient or RefusalDebt.
	Refusal Refusal
}

// Accepted reports whether o describes a successful direct acquisition.
//
// Accepted checks only the refusal classification. Use IsValid when callers
// also need to verify that the snapshot and refusal are structurally valid.
func (o Observation) Accepted() bool {
	return o.Refusal == RefusalNone
}

// Denied reports whether o describes ordinary bulkhead back-pressure.
//
// Denied is true for valid non-none capacity refusals. It does not classify
// programmer errors such as invalid amounts; those remain panics from the
// underlying capacity validation.
func (o Observation) Denied() bool {
	return o.Refusal.Refused()
}

// IsValid reports whether o is internally consistent.
//
// A valid observation requires a committed snapshot revision, a valid capacity
// snapshot value, and a valid refusal. Zero Observation is invalid because it
// does not describe a committed acquisition attempt.
func (o Observation) IsValid() bool {
	return !o.Snapshot.IsZeroRevision() &&
		o.Snapshot.Value.IsValid() &&
		o.Refusal.IsValid()
}
