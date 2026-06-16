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

// Package objectcache provides revisioned, queryable in-memory object cache
// primitives over objectstore.ListItem collections.
//
// Snapshot is immutable and represents a detached collection at one revision.
// Cache is mutable, concurrency-safe, and can replace state or apply committed
// objectstore changes. Both expose safe key lookup, full-item export, and
// objectquery-based filtering over cached state.
//
// The package owns materialized cached state and private indexes. Query
// semantics are delegated to api/objectquery: indexes only narrow candidates,
// and objectquery.Predicate.Match is always applied before returning query
// results.
//
// Cache preserves source list order. Apply appends creates in committed change
// order, preserves order for updates, and removes deleted keys while preserving
// the relative order of the remaining items. Apply assumes callers provide a
// complete ordered stream of committed changes for the cached collection; gap
// detection and recovery belong to future watch/reflector layers.
//
// The package does not read object stores, push filters into storage, parse
// selector strings, watch changes, run workload loops, own request policy, or
// validate resource descriptors. Indexes are intentionally private; there is no
// public objectindex package because indexes are useful only when tied to
// cache-owned state and change application.
package objectcache
