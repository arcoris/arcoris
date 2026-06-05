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

// ReserveResult is the status-rich result of Ledger.TryReserve.
type ReserveResult struct {
	// Status classifies the accounting result.
	Status ReserveStatus

	// Snapshot is the post-reservation state on success or unchanged current
	// state on refusal.
	Snapshot snapshot.Snapshot[Snapshot]

	// Reservation owns the demand on success. It is nil on refusal.
	Reservation *Reservation

	// Missing contains shortage diagnostics on insufficient or unknown-resource
	// refusal.
	Missing Vector

	// Debt contains debt diagnostics on debt refusal.
	Debt Vector
}

// Reserved reports whether the reserve attempt succeeded.
func (r ReserveResult) Reserved() bool {
	return r.Status.Reserved()
}

// Denied reports whether the reserve attempt was refused.
func (r ReserveResult) Denied() bool {
	return r.Status.Denied()
}
