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

// Package objectownership bridges object-level ownership state and ownership
// documents.
//
// Package fieldownership owns owner/path state, owner validation, conflict
// primitives, and deterministic field-level ownership ordering. Package
// fieldpath owns semantic path parsing and canonicalization. Package
// objectownership wraps those lower layers into object-surface ownership state
// and stable in-memory document shape. Package objectapply consumes State.
// Package objectstore stores Document. Package codecjson encodes and decodes
// Document. Package objectlifecycle decides when documents are stored,
// normalized, or committed.
//
// The package has three separate models. State is the operational object
// ownership state used by objectapply and lifecycle code. Document is stable,
// versioned, mutable in-memory document data used by stores and codecs. Normalize
// canonicalizes valid raw documents into deterministic document form. State is
// not a codec type. Document is not a wire format by itself.
//
// Validate accepts valid raw documents. Raw documents may contain duplicate
// owners, duplicate fields, unsorted entries, and empty entries. Normalize sorts
// owners, merges duplicate owners, deduplicates fields, sorts fields, prunes
// empty entries, and writes DocumentVersionV1. ValidateNormalized is available
// when callers need to require that a document is already canonical.
//
// Document version 1 owns only the Desired surface. Observed and metadata
// ownership are intentionally not modeled in v1; future document versions may
// add them explicitly.
//
// Document, Surface, and Entry are plain mutable document structs. They expose
// slices because they are raw document data. Use Clone, EntriesCopy, and
// FieldsCopy when caller-side detachment is required. Validate and Normalize do
// not make input immutable.
//
// The package does not apply objects, merge values, validate resources,
// validate descriptors, mutate metadata, access storage, assign storage
// revisions, run admission, authorize request subjects, execute lifecycle
// handlers, encode JSON/YAML/binary data, perform resource catalog lookup, or
// register runtime types. Runtime, storage, lifecycle, and codec layers own
// those concerns.
package objectownership
