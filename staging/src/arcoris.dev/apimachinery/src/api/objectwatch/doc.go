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
// A watch stream is scoped by Request.Collection: an objectstore.ListRequest
// containing resource plus structural namespace scope. The stream carries
// committed objectstore.Change values for that collection; it does not execute
// List, define objectquery filtering, parse selectors, read stores, store
// changes, apply cache mutations, run reflectors, or expose a transport.
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
// progress event revisions are monotonic non-decreasing boundaries.
//
// EventChanged carries one committed objectstore.Change. EventProgress is a
// watch progress marker, not a pagination bookmark or page token; it must not
// be applied to a cache and must never replace change delivery.
// EventRestartRequired is terminal: consumers must stop trusting the stream,
// perform a fresh list, and start a new watch from the new list revision.
//
// Stream is pull-based. Stream.Close must be idempotent. A source-side normal
// EOF is not part of the contract; a watch stream ends because the caller
// cancels or closes it, because continuity is explicitly lost, or because a
// restart is required.
//
// Validator can enforce the local request contract for consumers and tests. It
// checks that changed events match the watched objectstore.ListRequest resource
// and namespace scope, enforces Request.AllowProgress for progress events, and
// fails closed after continuity loss. Validator is not safe for concurrent use.
//
// Unsupported source capability is distinct from malformed input. Invalid
// request/start values return ErrInvalidRequest and ErrInvalidStart, while
// valid requests that a source cannot serve return ErrUnsupportedCapability.
//
// Consumers that need filtered watch semantics should combine this package with
// api/objectquery outside objectwatch, for example by projecting
// EventChanged.Change through objectquery.Predicate.ProjectChange. Future
// objectstorewatch, cache, reflector, runtime, transport, authorization,
// admission, and controller layers should depend on these contracts without
// moving their behavior into this package.
package objectwatch
