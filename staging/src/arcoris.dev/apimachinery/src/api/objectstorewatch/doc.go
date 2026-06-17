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

// Package objectstorewatch defines contracts that bridge objectstore
// collection reads and objectwatch streams.
//
// objectstore defines committed object state and structural collection
// matching. objectwatch defines how committed changes are streamed.
// objectstorewatch defines the list-to-watch boundary between them: a
// CollectionRead is a validated objectstore.ListResult tied to the exact
// objectstore.ListRequest that produced it, and a Boundary is the
// collection/revision pair from which a watch can continue.
//
// Package objectstorewatch deliberately avoids Snapshot terminology. In this
// repository, snapshots are point-in-time component read views provided by
// arcoris.dev/snapshot and related component packages. objectstorewatch works
// with validated collection reads used as list-to-watch boundaries. A
// CollectionRead is not a component snapshot and does not imply MVCC snapshot
// isolation unless a concrete implementation documents that stronger behavior.
//
// The central continuity contract is intentionally strict. For a CollectionRead S
// returned for collection C with boundary revision R, a watch request built
// from S.Boundary() must observe the same collection C. If Watch succeeds and the
// stream does not report restart-required or terminal continuity loss, it must
// deliver every matching committed change with revision greater than R in
// strictly increasing revision order. Silent gaps are forbidden.
//
// CollectionRead and Boundary are value-level contracts. They validate shape and
// preserve revision boundaries, but they do not prove that a concrete source
// still retains the required history or that a later watch will succeed. A
// concrete ListerWatcher proves continuity when it serves ListCollection and
// Watch consistently. If history is unavailable or continuity is lost,
// implementations must report that explicitly through objectwatch errors or
// restart-required events.
//
// This package does not implement storage, watch hubs, change logs,
// goroutines, caches, reflectors, transports, objectquery filtering,
// authorization, admission, lifecycle behavior, or resource descriptor
// validation. It also deliberately does not define mutation hooks such as
// OnCreate, OnUpdate, OnDelete, recorder, publisher, emitter, observer,
// journal, or change-log interfaces. Committed mutation capture is an
// implementation responsibility of future observable store wrappers or
// backend-native watch sources.
//
// Future implementations may use a write-through objectstore wrapper,
// backend-native watch support, a transactional outbox for persistent stores,
// or WAL/log tailing. Direct mutation of a wrapped backend outside such an
// implementation can break list-to-watch continuity.
package objectstorewatch
