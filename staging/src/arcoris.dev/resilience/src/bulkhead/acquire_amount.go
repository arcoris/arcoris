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

// TryAcquireAmount attempts to reserve amount in-flight capacity units.
//
// TryAcquireAmount is non-blocking and has the same ownership model as
// TryAcquire: success returns a live Lease that owns exactly amount units until
// Release or TryRelease returns them, while capacity exhaustion returns nil, an
// observed capacity state, and false. Observation.Refusal classifies the direct
// accounting outcome. The observation is produced through the underlying
// capacity ledger's explicit observed raw-reserve path. Under concurrent
// acquire, release, or limit changes its snapshot remains diagnostic rather
// than an exclusive serialization point.
//
// The successful path allocates the returned Lease token and uses raw ledger
// accounting underneath. It intentionally avoids creating a separate
// capacity.Reservation object for the same ownership event.
//
// The method intentionally does not wait, queue, apply fairness, observe context
// cancellation, retry, or classify denial as an error.
//
// Invalid amounts are programming or configuration errors, not denied
// acquisition results. Validation is delegated to the underlying capacity ledger
// after the Bulkhead receiver has been validated, preserving the package's
// existing panic behavior for nil or uninitialized receivers.
func (b *Bulkhead) TryAcquireAmount(amount Amount) (*Lease, Observation, bool) {
	b.requireReady()

	observation, ok := b.ledger.TryReserveObserved(amount)
	if !ok {
		return nil, Observation(observation), false
	}

	return &Lease{
		ledger: b.ledger,
		amount: amount,
	}, Observation(observation), true
}
