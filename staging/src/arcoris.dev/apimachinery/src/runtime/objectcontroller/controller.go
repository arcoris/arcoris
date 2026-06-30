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
	"sync"

	"arcoris.dev/apimachinery/runtime/objectreconciler"
)

// Controller owns worker lifecycle for object queue reconciliation.
//
// A Controller must not be copied after first use.
type Controller struct {
	// noCopy lets go vet report accidental Controller copies after first use.
	noCopy noCopy

	// queue provides work items and completion tracking.
	queue Queue

	// source provides the object snapshot consumed by each reconciliation
	// attempt.
	source objectreconciler.SnapshotSource

	// reconciler performs user reconciliation logic.
	reconciler objectreconciler.Reconciler

	// workers is the fixed worker count started by Run.
	workers int

	// mu guards running.
	mu sync.Mutex

	// running records whether Run is currently active.
	running bool
}
