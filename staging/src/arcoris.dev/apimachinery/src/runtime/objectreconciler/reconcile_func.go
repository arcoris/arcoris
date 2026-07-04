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

// ReconcileFunc adapts a function to Reconciler.
//
// A nil ReconcileFunc is treated as a wiring error and returns
// ErrNilReconciler instead of panicking.
type ReconcileFunc func(context.Context, Request, Snapshot) Result

// Reconcile calls f with ctx, request, and snap.
//
// Reconcile returns ErrNilReconciler when f is nil. It otherwise returns the
// function result as-is and does not recover panics.
func (f ReconcileFunc) Reconcile(ctx context.Context, request Request, snap Snapshot) Result {
	if f == nil {
		return Failure(ErrNilReconciler)
	}

	return f(ctx, request, snap)
}

// isNilReconciler reports whether reconciler is absent, including the typed nil
// ReconcileFunc case that a plain interface nil check cannot see.
func isNilReconciler(reconciler Reconciler) bool {
	if reconciler == nil {
		return true
	}
	fn, ok := reconciler.(ReconcileFunc)
	return ok && fn == nil
}
