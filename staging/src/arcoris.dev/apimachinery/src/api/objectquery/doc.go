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

// Package objectquery defines a stable query algebra for API object list items.
//
// A Query is an immutable expression value. The zero Query is valid and means
// All. Queries are constructed with boolean combinators and typed term
// constructors for storage keys, resources, labels, annotations, and explicitly
// registered selectable fields.
//
// Compile validates and canonicalizes a Query into a Predicate. A Predicate is
// immutable and concurrency-safe. It can match one objectstore.ListItem, filter
// already-loaded item slices while preserving input order, expose a detached
// canonical Query, provide conservative planning constraints, and project a
// committed objectstore.Change through the predicate.
//
// Query planning is advisory only. Plans may narrow cache or storage candidate
// sets, but Predicate.Match remains the final semantic source of truth.
// Negative requirements, NOT, and most OR expressions are intentionally kept
// residual unless they can be represented safely without changing semantics.
//
// Selectable fields are not arbitrary JSONPath or descriptor traversal.
// Resource-specific callers must register each queryable FieldRef through a
// SelectableFieldSet. Field paths use api/fieldpath.Path, whose elements model
// semantic field, map-key, list-index, and associative-list selector steps.
// Field terms currently evaluate Desired and Observed surfaces only; metadata
// labels and annotations have dedicated metadata terms. Missing means the path
// is absent. Null means the path exists and stores value.Null.
//
// SelectableField.Index affects only planning hints. IndexNone suppresses field
// constraints, IndexEquality exposes presence/equality/membership constraints,
// and IndexRange additionally exposes ordering constraints. Predicate.Match is
// always applied after any caller uses planning constraints, so planning cannot
// change query semantics.
//
// This package does not parse selector strings, decode HTTP query parameters,
// execute storage queries, build indexes, read stores, watch changes, run
// controllers, authorize callers, run admission, or validate resource
// descriptors. Those responsibilities belong to future adapters and higher
// layers.
package objectquery
