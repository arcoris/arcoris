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

// Snapshot returns the current revisioned capacity state.
//
// The returned value is the underlying capacity.Ledger snapshot. It is safe to
// store or compare as a value. It describes local in-flight capacity only; it
// does not include health, routing, scheduling, metrics, or distributed state.
func (b *Bulkhead) Snapshot() snapshot.Snapshot[Snapshot] {
	b.requireReady()
	return b.ledger.Snapshot()
}

// Revision returns the latest committed bulkhead capacity revision.
//
// Revisions are source-local to this Bulkhead. They are useful for cheap change
// detection by consumers observing the same Bulkhead, but they are not a global
// ordering across components.
func (b *Bulkhead) Revision() snapshot.Revision {
	b.requireReady()
	return b.ledger.Revision()
}
