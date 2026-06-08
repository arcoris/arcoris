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

package liveconfig

import (
	"sync"

	"arcoris.dev/snapshot"
)

// Holder owns the current last-good live configuration for one component or
// policy domain.
//
// A Holder serializes write-side Apply calls with mu, publishes accepted values
// through a snapshot.Publisher, and exposes Snapshot, Stamped, and Revision as
// read-side methods. Reads are delegated to the publisher and do not mutate
// holder state or execute validation, normalization, or source reload logic.
// This keeps read paths cheap and makes the holder suitable for read-mostly
// runtime policy such as limits, thresholds, retry knobs, and schedules.
//
// The holder does not own any input source. File watchers, environment readers,
// remote control-plane clients, and subscriber notification loops should call
// Apply from outside the package after they have built a candidate value.
//
// Holder is safe for concurrent use. Holder must be constructed with New; the
// zero value is invalid because it has no publisher or initial last-good value.
// Holder must not be copied after first use.
type Holder[T any] struct {
	// noCopy lets go vet report accidental Holder copies after first use.
	noCopy noCopy

	// mu serializes write-side Apply calls and protects lastErr.
	//
	// The snapshot.Publisher is safe for concurrent publishing, but Holder still
	// serializes Apply so the candidate preparation, equality check, publication,
	// and LastError update are observed as one coherent write-side operation.
	mu sync.Mutex

	// cfg contains the immutable construction policy used by New and Apply.
	//
	// Keeping these functions together makes the candidate pipeline explicit:
	// clone, then normalize, then validate, then optionally compare.
	cfg config[T]

	// pub publishes accepted immutable configuration values for read-side
	// Snapshot, Stamped, and Revision calls.
	pub *snapshot.Publisher[T]

	// lastErr is the last rejected Apply error. Successful changed and no-op
	// applies clear it.
	lastErr error
}
