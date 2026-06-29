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

// Package objectworkqueue provides a bounded, deduplicating runtime work queue
// for object-keyed reconciliation work.
//
// # Model
//
// A Queue stores pending object work between event sources, controllers, and
// future worker goroutines. Items are identified by objectstore.Key. Duplicate
// Add calls coalesce into one tracked item, and Get hands distinct queued items
// to callers in FIFO order.
//
// Queue tracks in-flight processing. If Add is called for an item while that
// item is processing, the item becomes dirty. Done requeues dirty processing
// items unless the queue is shutting down. This preserves the central invariant
// that work observed during processing is not lost.
//
// Capacity bounds distinct tracked items, not Add calls. Tracked items are the
// sum of queued and processing items.
//
// # Boundaries
//
// This package owns only queue state, deduplication, blocking handoff, dirty
// requeue semantics, and shutdown coordination. It does not execute reconcilers,
// read caches, watch stores, write objects, retry failures, implement backoff,
// schedule tasks, enforce task groups, perform capacity fitting, emit metrics,
// trace work, log events, or own controller lifecycle.
package objectworkqueue
