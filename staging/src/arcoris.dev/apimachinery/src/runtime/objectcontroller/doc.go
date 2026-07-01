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

// Package objectcontroller provides the worker lifecycle core that connects
// object work queues to object reconcilers.
//
// A Controller starts a fixed number of workers. Each worker takes one item
// from a queue, executes one read-only reconciliation attempt through
// objectreconciler.ReconcileOnce, and then attempts Queue.Done exactly once for
// the item it received. The controller waits for all workers before Run
// returns.
//
// The package sits after objectworkqueue and before objectreconciler in the
// runtime chain. It does not own list/watch synchronization, cache snapshots,
// queue producers, or queue shutdown. Queue shutdown is observed as clean
// completion when workers see objectworkqueue.ErrShutDown.
//
// Run stops when its context is cancelled, when the queue shuts down and
// drains, or when the first worker returns a fatal error. Fatal errors cancel
// the worker context so the remaining workers can exit, and Run still waits
// for every worker before returning.
//
// The controller does not recover panics from user reconciliation. After a
// successful Queue.Get, Queue.Done is attempted from a deferred path before a
// reconciliation panic continues through the worker goroutine. Because
// reconciliation is executed inside worker goroutines during Run, a panic is
// not converted into a Run return value. Panic recovery, panic-to-error
// conversion, logging, and retry policy belong outside this package.
//
// This package intentionally does not implement retry, backoff, rate limiting,
// delayed queues, scheduler policy, task groups, capacity fitting, object
// writes, finalizers, admission, metrics, tracing, logging, leader election,
// reflector lifecycle, event handler registries, panic policy hooks, or
// controller-runtime-style request mapping.
package objectcontroller
