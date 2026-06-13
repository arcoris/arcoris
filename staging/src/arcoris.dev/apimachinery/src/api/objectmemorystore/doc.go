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

// Package objectmemorystore provides an in-memory implementation of
// arcoris.dev/apimachinery/api/objectstore.Store.
//
// Store uses a fixed sharded key index. Shard locks protect only map structure:
// finding or creating the per-object slot. Each slot publishes immutable records
// through an atomic pointer, and live/update/delete transitions use
// compare-and-swap on that per-object pointer.
//
// Records are immutable after publication. Create and Update publish live
// records with normalized ownership and store-assigned revisions. Delete
// publishes a tombstone record instead of physically removing the slot.
// DeleteResult exposes both the deleted live state and the tombstone commit
// revision. Revision numbers are monotonic within one store, but concurrent CAS
// races may create gaps.
//
// The implementation is concurrency-safe for independent Store operations. It
// detaches caller input before publication and returns detached states from Get,
// List, Create, Update, and Delete.
//
// List supports live collection reads by resource and structural scope. It
// scans shards, copies matching slots, and atomically loads current records. It
// returns a detached current live collection read, not a historical MVCC
// snapshot under concurrent writes.
//
// The implementation is not durable, persistent, distributed, watch-capable,
// secondary-indexed, admission-aware, codec-aware, or transactional across
// keys. It does not validate resource descriptors, apply objects, stamp
// metadata, or execute lifecycle hooks.
package objectmemorystore
