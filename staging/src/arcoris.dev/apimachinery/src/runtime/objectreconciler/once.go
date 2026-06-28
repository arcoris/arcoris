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

package objectreconciler

import "context"

// ReconcileOnce executes one synchronous reconciliation attempt.
//
// ReconcileOnce reads exactly one snapshot from source, converts it to this
// package's Snapshot type, and calls reconciler at most once. It preserves
// source and user errors as-is and does not recover panics from user code.
// Source read errors, including objectcache.ErrNotReady, are returned without
// wrapping through Result.Err.
//
// ReconcileOnce checks ctx before reading the snapshot and again immediately
// before calling the reconciler. If either check observes cancellation,
// ReconcileOnce returns ctx.Err and does not call later steps.
//
// ReconcileOnce panics when ctx is nil. Use context.Background when no
// cancellation boundary is available.
func ReconcileOnce(ctx context.Context, source SnapshotSource, reconciler Reconciler) Result {
	if ctx == nil {
		panic("objectreconciler: nil context")
	}
	if err := ctx.Err(); err != nil {
		return Failure(err)
	}
	if source == nil {
		return Failure(ErrNilSource)
	}
	if isNilReconciler(reconciler) {
		return Failure(ErrNilReconciler)
	}

	sourceSnapshot, err := source.ReadSnapshot()
	if err != nil {
		return Failure(err)
	}
	if err := ctx.Err(); err != nil {
		return Failure(err)
	}

	return reconciler.Reconcile(ctx, snapshotFromSource(sourceSnapshot))
}
