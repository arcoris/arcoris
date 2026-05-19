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

package bulkhead

import "arcoris.dev/snapshot"

// TryAcquire attempts to reserve one in-flight slot without waiting.
//
// TryAcquire is the only admission operation in this package. It does not block,
// create waiters, observe context cancellation, retry, classify operation
// errors, or schedule work. If capacity is available, the returned Lease owns one
// capacity unit until Release or TryRelease returns it. If capacity is not
// available, TryAcquire returns nil, the current snapshot, and false.
//
// Rejected acquisition is not an error in this layer. It is ordinary bulkhead
// back-pressure and is represented by ok=false plus the observed snapshot. The
// rejected branch does not create a Lease because there is no capacity ownership
// to release.
func (b *Bulkhead) TryAcquire() (*Lease, snapshot.Snapshot[Snapshot], bool) {
	b.requireReady()

	reservation, snap, ok := b.ledger.TryReserve(1)
	if !ok {
		return nil, snap, false
	}

	return &Lease{reservation: reservation}, snap, true
}
