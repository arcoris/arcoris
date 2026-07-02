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

// Package objectenqueue maps committed object changes to object-keyed
// reconciliation work.
//
// The package is a small producer-side bridge between api/objectstore.Change
// and runtime/objectworkqueue.Item. Query-based filtering is delegated to
// api/objectquery.Predicate.ProjectChange. objectenqueue does not consume watch
// streams, mutate caches, run controller loops, reconcile objects, write
// objects, retry failures, rate limit producers, delay work, or own queue
// shutdown.
//
// objectstore.Change remains the value-level transition contract for committed
// store changes. The watch contract package owns stream continuity.
// objectworkqueue remains responsible for bounded storage, deduplication,
// blocking Add behavior, processing state, and dirty requeue semantics.
// objectenqueue only maps changes and forwards the mapped items to an Enqueuer.
package objectenqueue
