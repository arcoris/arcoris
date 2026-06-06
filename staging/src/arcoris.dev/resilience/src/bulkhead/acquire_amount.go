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

// TryAcquireAmount attempts to reserve amount in-flight capacity units.
//
// TryAcquireAmount is non-blocking and has the same ownership model as
// TryAcquire: success returns a live Lease that owns exactly amount units until
// Release or TryRelease returns them, while capacity exhaustion returns nil, an
// observed snapshot, and false. The snapshot is read after the accounting
// attempt. Under concurrent acquire, release, or limit changes it is a
// diagnostic observation, not an exclusive serialization point.
//
// The method intentionally does not wait, queue, apply fairness, observe context
// cancellation, retry, or classify denial as an error.
//
// Invalid amounts are programming or configuration errors, not denied admission
// decisions. Validation is delegated to the underlying capacity ledger after the
// Bulkhead receiver has been validated, preserving the package's existing panic
// behavior for nil or uninitialized receivers.
func (b *Bulkhead) TryAcquireAmount(amount Amount) (*Lease, snapshot.Snapshot[Snapshot], bool) {
	b.requireReady()

	ok := b.ledger.TryReserve(amount)
	snap := b.ledger.Snapshot()
	if !ok {
		return nil, snap, false
	}

	return &Lease{
		ledger: b.ledger,
		amount: amount,
	}, snap, true
}
