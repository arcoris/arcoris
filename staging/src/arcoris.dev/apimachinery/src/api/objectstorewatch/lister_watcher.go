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

package objectstorewatch

import "arcoris.dev/apimachinery/api/objectwatch"

// ListerWatcher can snapshot a collection and watch from snapshot boundaries.
//
// For a Snapshot S returned by Snapshot(ctx, C), a watch request built from
// S.Boundary() must observe the same collection C. If Watch succeeds and the
// stream does not report EventRestartRequired or a terminal continuity error,
// it must deliver every matching committed change with revision greater than
// S.Revision() in strictly increasing revision order.
//
// Implementations must not silently skip matching committed changes between
// Snapshot and Watch. If history is unavailable or continuity cannot be proven,
// the implementation must report that explicitly through objectwatch errors or
// restart-required events.
//
// Reflectors should usually depend on this interface rather than Store because
// reflectors do not need write methods.
type ListerWatcher interface {
	Snapshotter
	objectwatch.Source
}
