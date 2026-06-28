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

// Package objectreconciler defines read-only reconciliation primitives over
// runtime object cache snapshots.
//
// # Model
//
// One reconciliation attempt consumes exactly one detached objectcache.View at
// exactly one objectstore.Revision. ReconcileOnce reads a snapshot from a
// SnapshotSource, converts it into this package's Snapshot value, calls a user
// Reconciler synchronously, and returns the user's Result.
//
// The intended runtime chain is:
//
//	objectreflector -> objectcache -> objectreconciler -> future objectcontroller
//
// # Boundaries
//
// This package is intentionally read-only. User reconcilers should be
// idempotent and context-aware, but this package does not enforce idempotency,
// recover panics, write objects, patch status, schedule work, retry failures, or
// run background loops.
//
// # Non-goals
//
// This package does not implement watches, queues, workers, controllers, retry
// policy, backoff, metrics, tracing, panic recovery, event handlers, lifecycle,
// storage mutation, admission, validation, finalizers, scheduling policy, or
// leader election. Future objectworkqueue and objectcontroller packages may
// decide when and how often reconciliation attempts run.
package objectreconciler
