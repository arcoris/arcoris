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

import "errors"

var (
	// ErrNilSource reports that ReconcileOnce was called without a snapshot
	// source.
	ErrNilSource = errors.New("objectreconciler: nil snapshot source")

	// ErrNilReconciler reports that reconciliation logic was not provided.
	//
	// ReconcileFunc also returns ErrNilReconciler when the function value is
	// nil.
	ErrNilReconciler = errors.New("objectreconciler: nil reconciler")
)
