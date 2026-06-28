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

// Reconciler performs one user-defined reconciliation attempt.
//
// Reconciler implementations should be idempotent and context-aware. This
// package does not enforce idempotency and does not recover panics from
// Reconcile.
type Reconciler interface {
	// Reconcile executes one read-only attempt over snap and returns its
	// result.
	//
	// Reconcile receives a stable objectcache view at snap.Revision. It should
	// observe ctx for cancellation, but this package does not interrupt a
	// running call.
	Reconcile(context.Context, Snapshot) Result
}
