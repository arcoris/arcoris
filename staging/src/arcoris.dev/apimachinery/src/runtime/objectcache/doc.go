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

// Package objectcache provides a collection-bound runtime read-model cache for
// value-backed API objects.
//
// # Model
//
// Cache implements the reflector sink contract. A reflector installs complete
// collection reads through Replace and then applies committed objectstore.Change
// values through ApplyChange. The cache owns detached materialized latest state
// and, when configured, bounded per-object history of observed live versions and
// tombstones.
//
// # Boundaries
//
// This package is a concrete runtime implementation, not an API contract. It
// does not read from object stores, watch sources, run goroutines, perform
// retries, execute controllers, own query indexes, mutate storage, run
// admission, or serve transport APIs.
//
// # Query
//
// Query evaluates a compiled objectquery.Predicate over the latest live
// collection only. It scans cache-owned items in deterministic order, uses
// objectquery.Predicate.Match as the semantic source of truth, and returns
// detached matching items at the current cache revision. Historical object
// records are not considered by Query.
//
// Query indexes are intentionally deferred. Future private indexes may use
// objectquery planning hints to narrow candidates, but they must still confirm
// every result with Predicate.Match and must not change semantics.
//
// # Snapshots
//
// ReadSnapshot returns a detached View wrapped in
// snapshot.Snapshot[objectstore.Revision, View]. The snapshot revision is the
// objectstore collection revision, and it is the same revision returned by
// View.Revision.
//
// View provides stable Get, List, and Query reads over the latest live
// collection at one revision. It does not include per-object history records,
// tombstones, or future cache mutations. Historical point reads remain Cache
// methods through GetAt and PreviousLive; View does not reconstruct historical
// collection states.
//
// # History
//
// Historical reads are best-effort inside the configured per-object retention
// window. Replace resets retained history because a replacement may follow a
// continuity gap that the cache cannot repair. If a requested version is older
// than retained object history or crosses a replacement boundary, the cache
// reports ErrHistoryUnavailable instead of fabricating a result. Revisions are
// not assumed to be dense, so PreviousLive searches retained versions and never
// computes revision-1.
//
// Deleted records are bounded per key by RetainedVersionsPerObject, but the
// number of recently deleted keys can grow with deleted object cardinality. A
// future policy may add a global deleted-record retention limit if needed.
package objectcache
