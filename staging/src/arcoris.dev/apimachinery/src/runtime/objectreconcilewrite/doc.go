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

// Package objectreconcilewrite helps reconcilers build lifecycle write
// requests from the object currently visible in a reconciliation snapshot.
//
// The package does not execute writes. It only derives Resource, Object, and
// Expected fields for api/objectlifecycle write requests from the current
// object state selected by objectreconciler.Request and
// objectreconciler.Snapshot.
//
// Expected revision is always the current objectstore.State.Revision. The
// snapshot collection revision is deliberately not used as Expected because it
// describes the read-model boundary, not the committed revision of one object.
//
// objectreconcilewrite does not retry conflicts, requeue work, reread stores,
// start goroutines, patch status automatically, change the objectreconciler API,
// or own controller policy.
package objectreconcilewrite
