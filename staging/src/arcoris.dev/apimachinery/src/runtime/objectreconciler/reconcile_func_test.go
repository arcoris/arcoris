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

import (
	"context"
	"errors"
	"testing"
)

func TestNilReconcileFuncReturnsNilReconcilerError(t *testing.T) {
	var fn ReconcileFunc

	result := fn.Reconcile(context.Background(), Snapshot{})
	requireErrorIs(t, result.Err, ErrNilReconciler)
}

func TestReconcileFuncReceivesContextAndSnapshot(t *testing.T) {
	ctx := context.WithValue(context.Background(), testContextKey{}, "value")
	snap := snapshotFromSource(readySourceSnapshot(t, 4))
	var gotCtx context.Context
	var gotSnap Snapshot

	result := ReconcileFunc(func(ctx context.Context, snap Snapshot) Result {
		gotCtx = ctx
		gotSnap = snap
		return Success()
	}).Reconcile(ctx, snap)

	if result.Failed() {
		t.Fatalf("result = %#v; want success", result)
	}
	if gotCtx != ctx {
		t.Fatalf("ctx = %#v; want original context", gotCtx)
	}
	if gotSnap.Revision != snap.Revision || gotSnap.View.Revision() != snap.View.Revision() {
		t.Fatalf("snapshot = %#v; want %#v", gotSnap, snap)
	}
}

func TestReconcileFuncReturnsResultAsIs(t *testing.T) {
	want := Failure(errors.New("failed"))

	got := ReconcileFunc(func(context.Context, Snapshot) Result {
		return want
	}).Reconcile(context.Background(), Snapshot{})

	if got.Err != want.Err {
		t.Fatalf("result error = %v; want %v", got.Err, want.Err)
	}
}

type testContextKey struct{}
