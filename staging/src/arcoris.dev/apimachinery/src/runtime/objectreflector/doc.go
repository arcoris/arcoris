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

// Package objectreflector provides an active runtime synchronizer from an
// objectstorewatch.ListerWatcher source into a Sink.
//
// # Model
//
// A Reflector owns exactly one structural object collection. It reads that
// collection through ListCollection, installs the returned collection read into
// the sink, starts watching from the collection read boundary, and applies
// committed objectstore.Change values one at a time.
//
// # Continuity
//
// Request-aware stream validation is delegated to api/objectwatch.Validator.
// The Reflector owns source-to-sink control flow and Sink mutation, not the
// low-level watch event ordering rules. Watch continuity loss is handled by
// relisting; the Reflector never ignores a source gap, duplicate revision,
// out-of-collection change, or restart-required event in order to keep a stream
// running.
//
// RelistPolicy controls only the pacing between list-watch cycles after an
// explicit continuity loss, unavailable history, or restart-required event. It
// is not a retry framework and is not used for sink failures, malformed events,
// source contract violations, or invalid collection reads.
//
// # Sink Contract
//
// Sink receives committed changes only. Progress and restart events are
// stream-control signals handled inside the Reflector and are never forwarded as
// mutations. A failed sink operation is fatal because this package does not have
// an idempotency or repair contract for partially-applied sink writes.
//
// # Non-goals
//
// This package does not implement a cache, queue, informer, controller, retry
// framework, lifecycle manager, task group, metrics hook, transport adapter,
// object query engine, scheduler, resync period, event-handler registry, or
// persistent recovery mechanism. This core implementation relists immediately
// on continuity loss and restart-required events; retry and backoff are
// intentionally left to a future runtime layer that can make policy decisions
// outside the core source-to-sink protocol.
package objectreflector
