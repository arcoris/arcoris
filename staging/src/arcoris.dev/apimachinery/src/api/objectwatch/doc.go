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

// Package objectwatch defines contracts for ordered streams of committed
// objectstore.Change values.
//
// A watch stream is scoped by objectstore.ListRequest: resource plus structural
// namespace scope. The stream carries committed objectstore.Change values for
// that collection; it does not define objectquery filtering, parse selectors,
// read stores, apply cache mutations, run reflectors, or expose a transport.
//
// The core continuity contract is strict. For StartAfterRevision(R), if
// Source.Watch succeeds and the stream does not return EventRestartRequired or
// a terminal continuity error, the stream must deliver every committed change
// for the requested collection with change.Revision > R in strictly increasing
// revision order. Silent gaps are forbidden. If the source cannot prove
// continuity because history was compacted, reset, unavailable, or lost, it
// must fail Watch with a typed history/continuity error or emit
// EventRestartRequired and terminate the stream.
//
// objectstore.Revision is a store/source-local progress token. It is not
// wall-clock time, a global cluster version, or portable across unrelated
// sources unless a source explicitly documents that property. Within one
// successful stream, changed event revisions are strictly increasing and
// bookmark revisions are monotonic non-decreasing progress boundaries.
//
// EventChanged carries one committed objectstore.Change. EventBookmark is
// progress only; it must not be applied to a cache and must never replace
// change delivery. EventRestartRequired is terminal: consumers must stop
// trusting the stream, perform a fresh list, and start a new watch from the new
// list revision.
//
// Consumers that need filtered watch semantics should combine this package with
// api/objectquery outside objectwatch, for example by projecting
// EventChanged.Change through objectquery.Predicate.ProjectChange. Future
// cache, reflector, runtime, transport, authorization, admission, and
// controller layers should depend on these contracts without moving their
// behavior into this package.
package objectwatch
