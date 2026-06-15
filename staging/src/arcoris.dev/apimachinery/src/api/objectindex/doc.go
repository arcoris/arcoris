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

// Package objectindex provides optional static in-memory indexes for
// already-loaded api/objectstore.ListItem collections.
//
// objectindex accelerates api/objectquery selection without defining query
// semantics. objectquery remains the source of truth: every indexed selection
// narrows candidates when possible and then applies the compiled objectquery
// predicate before returning final items. Results are therefore equivalent to a
// full objectquery predicate scan over the same input, including order and
// duplicate input items.
//
// Indexes are static and shallow. Build copies the input item slice structure
// but does not clone objectstore.State, object metadata maps, or payload values.
// Select returns shallow ListItem copies. Callers that need detached object
// states must build the index from detached items or clone returned items
// themselves.
//
// The package does not read object stores, push filters into storage, maintain
// live caches, watch changes, process events, own revision semantics, parse
// selector strings, validate resource descriptors, or implement serving,
// policy, workload-orchestration, pagination, or MVCC behavior.
package objectindex
