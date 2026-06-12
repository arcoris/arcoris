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

// Package objectlifecycle coordinates local API object lifecycle transitions
// over api/objectstore.
//
// # Package responsibility
//
// objectlifecycle is the descriptor-aware stateful orchestration layer for
// value-backed API objects. It resolves resource contracts, prepares lifecycle
// requests, delegates object and value semantics to lower packages, and commits
// already-computed state through objectstore.
//
// # Relationship to lower layers
//
// api/objectvalidation validates object envelopes against resolved resource
// contracts. api/objectapply computes pure Desired apply output for existing
// live objects. api/objectownership represents canonical object ownership
// state. api/objectstore commits already-computed state with optimistic
// concurrency. objectlifecycle composes these layers into Get,
// Create, Apply, UpdateObserved, PatchMetadata, and Delete operations.
//
// # Operation semantics
//
// Get resolves an explicit group/version/resource and object identity, reads
// committed live state, and returns either the found state or ErrNotFound. Get
// does not revalidate stored Desired or Observed payloads against current
// descriptors; descriptor-aware validation is a write-path responsibility so
// reads remain stable under descriptor evolution.
//
// Create accepts a value-backed live object envelope, resolves its resource by
// group/version/kind, validates it against the resource contract, initializes
// Desired ownership from explicitly present Desired fields, and commits through
// objectstore.Create. Create may accept Observed when the selected resource
// version defines an Observed surface because it commits complete live state; it
// is not an admission-layer user create request.
//
// Apply accepts Desired intent. It resolves the resource by group/version/kind,
// rejects Observed apply intent before checking whether the object exists,
// validates object shape, and reads live state. Missing-object Apply creates the
// object by initializing ownership with the same field extraction path as
// Create. Existing-object Apply delegates Desired merge and ownership semantics
// to api/objectapply and commits through objectstore.Update.
//
// UpdateObserved resolves an explicit group/version/resource and object
// identity, requires an existing live object and non-zero expected store
// revision, validates the replacement Observed payload against the selected
// resource version's Observed descriptor, replaces Observed, updates Observed
// ownership, preserves Desired and metadata ownership, and commits through
// objectstore.Update.
//
// PatchMetadata resolves an explicit group/version/resource and object identity,
// requires an existing live object and non-zero expected store revision, and
// patches labels and annotations only. Nil patch values delete keys and non-nil
// values set keys. The operation preserves TypeMeta, ObjectMeta identity/system
// fields, Desired, Observed, Desired ownership, and Observed ownership.
// Finalizers and ownerReferences are intentionally not generic metadata patch
// fields.
//
// Delete resolves an explicit group/version/resource and object identity,
// requires a non-zero expected store revision, and commits a tombstone through
// objectstore.Delete. The returned State is the deleted live state and keeps the
// previous live revision. Result.Revision is the tombstone commit revision.
//
// # Apply create-on-missing policy
//
// Apply creates missing objects by design. This is server-side-apply-like
// behavior. Existing-object Apply delegates to objectapply; missing-object Apply
// initializes ownership through valuefieldset via the same helper used by
// Create.
//
// # Observed policy
//
// ApplyRequest represents Desired apply intent. Applied Observed payloads are
// rejected even when the resource defines Observed. Live Observed is preserved
// during existing-object Apply.
//
// UpdateObserved is the separate Observed writer operation. It does not mutate
// Desired or metadata. The first implementation uses complete replacement
// semantics for Observed and assigns Observed ownership to the supplied owner
// for the fields explicitly present in the replacement payload.
//
// # Metadata policy
//
// PatchMetadata is limited to labels and annotations. It does not mutate
// metadata.name, metadata.namespace, metadata.uid, metadata.resourceVersion,
// metadata.generation, metadata.createdAt, metadata.deletion, finalizers, or
// ownerReferences. It does not stamp resourceVersion or generation.
//
// # Non-goals
//
// objectlifecycle does not decode wire formats, select codecs, serve HTTP or
// gRPC, run admission, authorize subjects, stamp ObjectMeta resourceVersion or
// generation, generate UIDs, execute finalizers, perform graceful deletion,
// default, convert, prune, list, watch, compact tombstones, reconcile
// controllers, emit metrics/logging/tracing, or start background goroutines.
package objectlifecycle
