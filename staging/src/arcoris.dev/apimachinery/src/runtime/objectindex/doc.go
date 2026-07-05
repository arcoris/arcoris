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

// Package objectindex maintains secondary indexes over reflected object
// collections.
//
// An Index is a passive read-side primitive fed through the same Sink methods as
// runtime/objectreflector consumers: Replace rebuilds every index from a list
// boundary, and ApplyChange incrementally updates memberships for committed
// changes. The package owns no source, store, cache, queue, controller, or
// reconciler, and starts no goroutines.
//
// Index definitions are caller supplied. Extractors map one listed object to
// zero or more opaque, normalized string values. The package stores those values
// under a named index and supports reverse lookups from name/value to
// objectstore keys.
//
// objectindex is useful for mapped and multi-source controller mappers that
// need fast reverse lookups, such as finding target object keys affected by a
// changed source object. It does not execute objectquery.Predicate, replace
// runtime/objectcache, plan queries, write objects, retry failures, or implement
// scheduling policy.
package objectindex
