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

// Package object defines generic ARCORIS API object and list envelopes.
//
// Package meta owns TypeMeta, ObjectMeta, PageMeta, and metadata validation.
// Package resource owns durable resource-family contracts. Package
// objectvalidation binds Object[D, O] values to resource definitions and typed
// surface validators. Package objectapply performs value-backed desired apply,
// and objectstore stores committed value-backed object state. Package object
// does not perform those responsibilities.
//
// Object[D, O] is:
//
//   - TypeMeta
//   - ObjectMeta
//   - Desired D
//   - optional Observed *O
//
// D is the desired/requested payload type. O is the observed/computed payload
// type. Observed is absent when nil.
//
// List[T] is:
//
//   - TypeMeta
//   - PageMeta
//   - Items []T
//
// T may be Object[D, O], a concrete resource type, a decoded adapter type, or
// another API object representation.
//
// Constructors clone metadata values using metadata clone semantics. NewList
// shallow-copies the item slice and preserves nil-vs-empty shape. Observed
// values are stored behind a fresh pointer when set through constructors and
// helpers. Desired, Observed, and Items payload values are not deep-copied.
// Exported fields remain mutable like ordinary Go struct fields; use helper
// methods when caller-side metadata or slice detachment is desired.
//
// ValidateMeta validates only TypeMeta and ObjectMeta for Object, and only
// TypeMeta and PageMeta for List. Desired, Observed, and Items are
// intentionally ignored. Resource-aware validation belongs to api/objectvalidation.
// Value-backed payload validation belongs to api/valuevalidation through
// objectvalidation or objectapply.
//
// Struct tags provide simple representation hints. Canonical wire formats and
// manifest codecs belong to codec packages, not api/object.
//
// Non-goals: no resource matching, descriptor-aware payload validation, apply,
// ownership, lifecycle, storage, watches, admission, defaulting, conversion,
// pruning, selector matching, codecs, runtime schemes, clients, controllers,
// catalogs, or global registration.
package object
