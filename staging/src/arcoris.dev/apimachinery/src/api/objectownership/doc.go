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

// Package objectownership defines canonical object-level ownership state.
//
// Package fieldownership owns owner/path state, owner validation, conflict
// primitives, and deterministic field-level ownership ordering. Package
// fieldpath owns semantic path parsing and canonicalization. Package
// objectsurface names the stable object surface taxonomy. Package
// objectownership wraps those lower layers into multi-surface object ownership
// state. Package objectapply consumes the Desired portion of State. Package
// objectstore stores State. Package codecjson encodes and decodes State. Package
// objectlifecycle decides when ownership state is initialized, updated, or
// committed.
//
// State is the single canonical ownership model. There is no separate ownership
// document model, document version, or migration layer in this package. Stores,
// codecs, lifecycle code, and apply code all exchange State directly.
//
// State stores ownership in separate surface-relative namespaces:
// Desired, Observed, metadata labels, and metadata annotations. A path such as
// $.ready in Observed is unrelated to $.ready in Desired, and a metadata label
// path such as $["scheduler.arcoris.dev/mode"] is relative to the labels
// surface rather than to a synthetic whole-object $.metadata path.
//
// Normalize canonicalizes each surface independently by relying on
// fieldownership's deterministic owner ordering, duplicate-owner merge,
// duplicate-field deduplication, and empty-entry pruning. Validate checks each
// surface with fieldownership structural validation. ValidateNormalized is
// available for callers that need to assert already-canonical state.
//
// The package does not apply objects, merge values, validate resources,
// validate descriptors, mutate metadata, access storage, assign storage
// revisions, run admission, authorize request subjects, execute lifecycle
// handlers, encode JSON/YAML/binary data, perform resource catalog lookup, or
// register runtime types. Runtime, storage, lifecycle, and codec layers own
// those concerns.
package objectownership
