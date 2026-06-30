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

package objectcontroller

import "arcoris.dev/apimachinery/runtime/objectreconciler"

// New validates dependencies and returns a Controller.
//
// New rejects invalid worker counts, missing queue consumers, missing snapshot
// sources, and missing reconcilers. It accepts only the queue consumer
// interface; queue producers and shutdown ownership remain outside this
// package.
func New(
	opts Options,
	queue Queue,
	source objectreconciler.SnapshotSource,
	reconciler objectreconciler.Reconciler,
) (*Controller, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	if queue == nil {
		return nil, ErrNilQueue
	}
	if source == nil {
		return nil, ErrNilSource
	}
	if isNilReconciler(reconciler) {
		return nil, ErrNilReconciler
	}

	return &Controller{
		queue:      queue,
		source:     source,
		reconciler: reconciler,
		workers:    opts.Workers,
	}, nil
}

// isNilReconciler reports whether reconciler is absent, including the typed
// nil ReconcileFunc case that a plain interface nil check cannot see.
func isNilReconciler(reconciler objectreconciler.Reconciler) bool {
	if reconciler == nil {
		return true
	}
	fn, ok := reconciler.(objectreconciler.ReconcileFunc)
	return ok && fn == nil
}
