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
// Store uses a fixed sharded key index and one atomic publication slot per
// object. Shard locks protect only map structure. Live object transitions use
// per-object compare-and-swap over immutable records, so unrelated objects do
// not serialize behind one global state mutex.
//
// Records are immutable after publication. Delete commits a tombstone record
// instead of physically removing a slot; this avoids delete/update races and
// leaves a clean foundation for future watch/event and compaction layers.
//
// The implementation is intended for tests, local prototypes, and future
// runtime composition. It is not durable, distributed, persistent, indexed,
// list-capable, or watch-capable. It does not validate objects against resource
// descriptors, apply objects, run admission, or encode/decode wire formats.
package objectmemorystore
