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

// Package objectenqueue maps reflected object state and committed object
// changes to object-keyed reconciliation work.
//
// The package is a small producer-side bridge from api/objectstore values to
// runtime/objectworkqueue.Item. Query-based filtering is delegated to
// api/objectquery.Predicate.Match and api/objectquery.Predicate.ProjectChange.
// objectenqueue does not consume watch streams, mutate read models, run control
// loops, reconcile objects, write objects, retry failures, rate limit
// producers, delay work, own queue shutdown, or implement scheduling policy.
//
// objectstore.Change remains the value-level transition contract for committed
// store changes. The watch contract package owns stream continuity.
// objectworkqueue remains responsible for bounded storage, deduplication,
// blocking Add behavior, processing state, and dirty requeue semantics.
// objectenqueue only maps state/change inputs and forwards the mapped items to
// an Enqueuer.
//
// ReflectorSink is stateful and replace-aware. It keeps a local known set only
// to avoid losing reconciliation work across replace/relist boundaries. That
// known set is not a read model and must not replace the runtime read-model
// layer.
package objectenqueue
