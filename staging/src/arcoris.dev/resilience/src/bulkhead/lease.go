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

import "arcoris.dev/capacity"

// Lease owns in-flight bulkhead capacity until it is released.
//
// Lease is the bulkhead-domain name for a capacity.Reservation. The underlying
// reservation enforces release ownership and updates the owning ledger. Lease
// adds no counters of its own and intentionally exposes only the operations that
// are meaningful for protected in-flight work.
//
// Lease must not be copied after creation.
type Lease struct {
	// noCopy lets go vet report accidental Lease copies after first use.
	noCopy noCopy

	// reservation owns the low-level capacity units.
	reservation *capacity.Reservation
}

// Amount returns the number of in-flight capacity units owned by l.
//
// TryAcquire creates one-unit leases, while TryAcquireAmount and TryAdmit may
// create weighted leases. The amount is immutable after acquisition and remains
// observable before and after release.
func (l *Lease) Amount() Amount {
	l.requireReady()
	return Amount(l.reservation.Amount())
}

// Released reports whether l has already returned its capacity to the Bulkhead.
//
// Released is an ownership-state query. It does not release capacity and does
// not make Release idempotent; callers that need idempotent cleanup should use
// TryRelease.
func (l *Lease) Released() bool {
	l.requireReady()
	return l.reservation.Released()
}
