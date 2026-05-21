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

import (
	"arcoris.dev/capacity"
	"arcoris.dev/snapshot"
)

// SetLimit replaces the in-flight limit and returns the resulting snapshot.
//
// Active leases are never revoked. Reducing the limit below active leases uses
// capacity debt semantics: Available becomes zero, Debt reports the excess
// active leases, and new acquisitions are denied until enough leases are
// released. This package does not add extra drain, eviction, or cancellation
// policy on top of capacity.Ledger.
//
// A zero limit is valid and closes the bulkhead until a later SetLimit raises
// capacity again.
func (b *Bulkhead) SetLimit(limit Amount) snapshot.Snapshot[Snapshot] {
	b.requireReady()
	return b.ledger.SetLimit(capacity.Amount(limit))
}
