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

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectreconciler"
)

func TestNewRejectsInvalidWorkers(t *testing.T) {
	_, err := New(Options{}, &recordingQueue{}, &fakeSnapshotSource{}, &fakeReconciler{})
	requireErrorIs(t, err, ErrInvalidWorkers)
}

func TestNewRejectsNilQueue(t *testing.T) {
	_, err := New(Options{Workers: 1}, nil, &fakeSnapshotSource{}, &fakeReconciler{})
	requireErrorIs(t, err, ErrNilQueue)
}

func TestNewRejectsNilSource(t *testing.T) {
	_, err := New(Options{Workers: 1}, &recordingQueue{}, nil, &fakeReconciler{})
	requireErrorIs(t, err, ErrNilSource)
}

func TestNewRejectsNilReconciler(t *testing.T) {
	_, err := New(Options{Workers: 1}, &recordingQueue{}, &fakeSnapshotSource{}, nil)
	requireErrorIs(t, err, ErrNilReconciler)
}

func TestNewRejectsNilReconcileFunc(t *testing.T) {
	var fn objectreconciler.ReconcileFunc
	_, err := New(Options{Workers: 1}, &recordingQueue{}, &fakeSnapshotSource{}, fn)
	requireErrorIs(t, err, ErrNilReconciler)
}

func TestNewAcceptsValidDependencies(t *testing.T) {
	controller, err := New(
		Options{Workers: 2},
		&recordingQueue{},
		&fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		objectreconciler.ReconcileFunc(func(context.Context, objectreconciler.Snapshot) objectreconciler.Result {
			return objectreconciler.Success()
		}),
	)
	requireNoError(t, err)
	if controller == nil {
		t.Fatal("controller is nil")
	}
}
