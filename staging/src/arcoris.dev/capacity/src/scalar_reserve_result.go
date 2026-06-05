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

// ScalarReserveResult is the status-rich result of ScalarLedger.TryReserve.
type ScalarReserveResult struct {
	// Status classifies the accounting result.
	Status ReserveStatus

	// Snapshot is the post-reservation state on success or unchanged current
	// state on refusal.
	Snapshot snapshot.Snapshot[ScalarSnapshot]

	// Reservation owns the amount on success. It is nil on refusal.
	Reservation *ScalarReservation
}

// Reserved reports whether the scalar reserve attempt succeeded.
func (r ScalarReserveResult) Reserved() bool {
	return r.Status.Reserved()
}

// Denied reports whether the scalar reserve attempt was refused.
func (r ScalarReserveResult) Denied() bool {
	return r.Status.Denied()
}
