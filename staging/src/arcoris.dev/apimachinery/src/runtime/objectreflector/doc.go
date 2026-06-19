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
// A Reflector lists one structural object collection, replaces the sink with
// that collection, starts watching from the collection read boundary, and
// applies committed changes one at a time. When a watch stream cannot preserve
// continuity, the Reflector relists instead of allowing silent gaps.
//
// Sink receives committed objectstore.Change values only. Progress and restart
// events are stream-control signals handled inside the Reflector and are never
// forwarded as mutations.
//
// The first implementation intentionally keeps retry and backoff outside this
// package. A failed sink operation is fatal because the Reflector does not know
// whether the sink partially applied the operation or how to repair it.
//
// This package does not implement a cache, queue, informer, controller, retry
// framework, lifecycle manager, task group, metrics hook, transport adapter,
// object query engine, scheduler, resync period, event-handler registry, or
// persistent recovery mechanism.
package objectreflector
