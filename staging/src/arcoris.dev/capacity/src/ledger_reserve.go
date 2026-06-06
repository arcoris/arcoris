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

package capacity

import "arcoris.dev/snapshot"

// TryReserve attempts to reserve amount using the raw scalar accounting path.
//
// TryReserve is the allocation-free hot path for callers that already own their
// own lifecycle token. It mutates only accounting counters and returns false for
// ordinary capacity refusal. A zero amount is a programmer error and panics.
func (l *Ledger) TryReserve(amount Amount) bool {
	return l.reserveAmount(amount)
}

// TryReserveObserved reserves amount through the raw path and returns diagnostics.
//
// Observation is built only because this observed method was explicitly called.
// Refusal describes this attempt. Snapshot is read after the attempt and is
// therefore exact for single-threaded callers and observational under concurrent
// mutation, matching Ledger's snapshot consistency contract.
func (l *Ledger) TryReserveObserved(amount Amount) (Observation, bool) {
	refusal, ok := l.tryReserveAmount(amount)

	return Observation{
		Snapshot: l.Snapshot(),
		Refusal:  refusal,
	}, ok
}

// TryAcquire reserves amount and returns a capacity-owned Reservation token.
//
// TryAcquire uses the same accounting path as TryReserve. On success, the only
// steady-state allocation is the returned Reservation object. On refusal, it
// returns nil, false without constructing diagnostics.
func (l *Ledger) TryAcquire(amount Amount) (*Reservation, bool) {
	if !l.reserveAmount(amount) {
		return nil, false
	}

	return &Reservation{ledger: l, amount: amount}, true
}

// TryAcquireObserved reserves amount and returns an owned token plus diagnostics.
func (l *Ledger) TryAcquireObserved(amount Amount) (*Reservation, Observation, bool) {
	observation, ok := l.TryReserveObserved(amount)
	if !ok {
		return nil, observation, false
	}

	return &Reservation{ledger: l, amount: amount}, observation, true
}

// Release returns raw accounting capacity to l.
//
// Release is strict. It panics if amount is zero or if the ledger does not
// currently have enough reserved capacity. Callers that need idempotent cleanup
// should use an owned Reservation or guard release in their own owner object.
func (l *Ledger) Release(amount Amount) {
	if !l.TryRelease(amount) {
		panicAt("ledger.reserved", ErrReservedUnderflow, "reserved amount is smaller than release amount")
	}
}

// TryRelease returns raw accounting capacity to l when possible.
//
// TryRelease is useful for owner objects that want idempotent cleanup while
// still avoiding a capacity Reservation allocation. A zero amount is a
// programmer error and panics.
func (l *Ledger) TryRelease(amount Amount) bool {
	return l.releaseAmount(amount)
}

// ReleaseObserved returns raw accounting capacity and then reads a snapshot.
func (l *Ledger) ReleaseObserved(amount Amount) snapshot.Snapshot[Snapshot] {
	l.Release(amount)

	return l.Snapshot()
}

// TryReleaseObserved returns raw accounting capacity, a snapshot, and the outcome.
func (l *Ledger) TryReleaseObserved(amount Amount) (snapshot.Snapshot[Snapshot], bool) {
	ok := l.TryRelease(amount)

	return l.Snapshot(), ok
}

// reserveAmount is the raw scalar reserve CAS loop.
func (l *Ledger) reserveAmount(amount Amount) bool {
	_, ok := l.tryReserveAmount(amount)

	return ok
}

// tryReserveAmount reserves amount or reports the accounting refusal.
func (l *Ledger) tryReserveAmount(amount Amount) (Refusal, bool) {
	requirePositiveAmount(amount)
	l.requireReady()

	for {
		limit := Amount(l.limit.Load())
		reserved := Amount(l.reserved.Load())
		refusal := refusalForScalar(limit, reserved, amount)
		if refusal.Refused() {
			return refusal, false
		}

		next, ok := reserved.CheckedAdd(amount)
		if !ok {
			return RefusalInsufficient, false
		}
		if l.reserved.CompareAndSwap(reserved.Uint64(), next.Uint64()) {
			l.advanceRevision()
			return RefusalNone, true
		}
	}
}

// releaseAmount is the raw scalar release CAS loop.
func (l *Ledger) releaseAmount(amount Amount) bool {
	requirePositiveAmount(amount)
	l.requireReady()

	for {
		reserved := Amount(l.reserved.Load())
		next, ok := reserved.CheckedSub(amount)
		if !ok {
			return false
		}
		if l.reserved.CompareAndSwap(reserved.Uint64(), next.Uint64()) {
			l.advanceRevision()
			return true
		}
	}
}

// refusalForScalar classifies a scalar accounting attempt.
func refusalForScalar(limit Amount, reserved Amount, amount Amount) Refusal {
	if reserved > limit {
		return RefusalDebt
	}
	if amount > limit-reserved {
		return RefusalInsufficient
	}

	return RefusalNone
}
