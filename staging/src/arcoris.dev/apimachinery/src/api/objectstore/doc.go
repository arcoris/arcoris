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

// Package objectstore defines the committed-state contract for value-backed API
// objects.
//
// State contains a value-backed live object envelope, object ownership state,
// and a store-local Revision. Revisions are assigned by a Store when Create,
// Update, or Delete commits. Committed State and Change values must carry
// non-zero revisions. List reports the store revision observed during the
// collection read; that list boundary may be zero for an empty store before any
// committed mutation has been observed. Revisions are not API resourceVersion
// values, ObjectMeta generations, wall-clock timestamps, durable storage
// versions, or globally comparable values across stores. They may have gaps.
//
// Callers are responsible for resolving resources, validating object envelopes,
// computing apply results, deciding admission, and preparing lifecycle
// semantics before calling a Store. Store implementations validate keys and
// input state shape, detach caller state, normalize ownership before commit,
// assign revisions, enforce optimistic concurrency, and return detached
// committed states.
//
// Input State values for Create and Update must have zero Revision. Committed
// State values returned by stores have non-zero Revision and normalized
// ownership state. ValidateInputState accepts valid ownership state;
// PrepareInputState clones and normalizes it. ValidateCommittedState requires
// normalized ownership.
//
// Delete commits a tombstone and returns DeleteResult. DeleteResult.Deleted is
// the live state that was deleted and keeps its previous live revision.
// DeleteResult.Revision is the store-local tombstone commit revision.
//
// Change describes one structurally valid committed transition: created,
// updated, or deleted. It is a reusable value model for future watch, cache,
// audit, or controller layers, but objectstore does not implement those layers
// and Change values do not prove that a transition occurred in a particular
// Store instance.
//
// List reads committed live states for one resource collection and structural
// scope. It returns only live records; missing, deleted, and tombstoned objects
// are omitted. ListResult.Revision is the store revision observed by the
// operation. Store implementations need not provide historical MVCC snapshot
// isolation unless a concrete implementation documents that stronger behavior.
// Structural collection matching helpers in this package check only resource
// equality and ListScope namespace membership.
//
// The package deliberately does not validate resource descriptors, apply
// objects, compute field conflicts, run admission or authorization, default,
// convert, prune, watch, reflect changes, provide List -> Watch continuity,
// encode/decode wire formats, expose serving behavior, apply selectors,
// paginate results, or stamp object metadata resourceVersion/generation fields.
// Those responsibilities belong to higher lifecycle, apply, resource, codec,
// serving, future objectstorewatch, and future watch/runtime layers.
//
// The in-memory implementation lives in
// arcoris.dev/apimachinery/runtime/objectmemorystore.
package objectstore
