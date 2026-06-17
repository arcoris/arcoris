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
// The package is an implementation of the api/objectstorewatch contracts, not a
// new storage engine and not a reflector. Store accepts any objectstore.Store,
// delegates committed state operations to that backend, records committed
// objectstore.Change values produced by Create, Update, and Delete, and serves
// objectwatch streams from bounded in-memory retained history plus live fanout.
//
// After wrapping a backend, callers must route all mutations through the
// wrapper. Direct backend mutation bypasses change capture, retained history,
// and live watcher dispatch, making ListCollection -> Watch continuity
// unverifiable. This is a caller contract; it is not generically detectable
// through the objectstore.Store interface.
//
// Store is passive. It creates no background dispatcher goroutine, no progress
// ticker, no reflector, no cache, no retry loop, no backoff policy, and no
// lifecycle controller. History is bounded and in-memory only; once compacted,
// old starts fail explicitly with objectwatch history errors. Slow watchers
// lose continuity instead of blocking writers, and matching events are never
// silently dropped as successful continuation.
package objectstorewatch
