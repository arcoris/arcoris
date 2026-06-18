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

// Package objectstorewatch provides a runtime write-through observable wrapper
// for objectstore.Store implementations.
//
// # Model
//
// The package is an implementation of the api/objectstorewatch contracts, not a
// new storage engine and not a reflector. Store accepts any objectstore.Store
// backend, delegates committed state operations to that backend, captures
// committed objectstore.Change values produced by Create, Update, and Delete,
// and serves objectwatch streams from retained history plus live fanout.
//
// The intended ownership chain is:
//
//	backend objectstore.Store
//	        -> runtime/objectstorewatch.Store
//	        -> objectstore.Store + CollectionLister + objectwatch.Source
//
// # Continuity
//
// Store serializes backend calls, ListCollection, Watch registration, retained
// history mutation, and live fanout with one mutex. That single serialization
// point is the continuity proof available over a generic objectstore.Store,
// which does not expose a transactional watch hook or current revision API.
// Matching committed changes are either delivered, retained for replay, or
// reported as explicit continuity loss. Silent gaps are not successful
// continuation.
//
// # History Retention
//
// Retained history is bounded and in-memory only. It is useful for short
// historical starts and reconnect windows, not durable recovery. Once the
// retained prefix is compacted, starts before the compaction boundary fail with
// objectwatch history-unavailable errors.
//
// # Backpressure
//
// Watch streams are pull-based. StreamBuffer bounds live backlog only;
// historical replay is stored separately so small live buffers do not make
// retained history unwatchable. Writers never wait for slow streams. A slow
// stream loses continuity with ErrStreamOverflow wrapped as objectwatch
// continuity loss, while other streams and future watchers may continue.
//
// # Backend Ownership
//
// After wrapping a backend, callers must route all mutations through the
// wrapper. Direct backend mutation bypasses change capture, retained history,
// and live stream dispatch, making ListCollection -> Watch continuity
// unverifiable. This is a caller contract; it is not generically detectable
// through the objectstore.Store interface.
//
// # Non-Goals
//
// Store is passive. It creates no background dispatcher goroutine, progress
// ticker, reflector, cache, retry loop, backoff policy, lifecycle controller,
// objectquery filter, transport adapter, persistent WAL, or transactional
// outbox. Those concerns belong in higher runtime layers.
package objectstorewatch
