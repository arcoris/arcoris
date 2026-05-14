/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package run

import "context"

// Task is a context-aware unit of runtime work owned by a Group.
//
// A Task is the smallest executable contract in package run. It represents one
// long-lived or short-lived unit of component runtime work, such as a dispatcher
// loop, controller loop, worker loop, sampler, watcher, drain routine, cleanup
// routine, or any other goroutine owned by a component runtime scope.
//
// The Task contract is intentionally context-first. Package run does not define
// a separate interrupt callback, signal callback, lifecycle callback, retry
// callback, logger hook, metrics hook, or tracing hook. The context passed to a
// Task is the only shutdown and cancellation channel understood by this package.
// Higher-level owners may derive that context from process signals, lifecycle
// transitions, explicit owner cancellation, tests, timeouts, or other runtime
// coordination mechanisms.
//
// Task implementations SHOULD observe ctx when they can block, wait, receive,
// send, poll, sleep, acquire resources, call external systems, or run loops that
// may outlive the immediate caller. Short pure in-memory work MAY ignore ctx when
// it cannot reasonably block and always finishes promptly.
//
// A Task returns nil when it completed successfully or stopped cleanly as part of
// expected owner-driven shutdown. A non-nil error means the task failed. Depending
// on the owning Group configuration, a non-nil task error may cancel the shared
// group context and cause sibling tasks to stop.
//
// Expected context-driven shutdown SHOULD usually return nil:
//
//	func(ctx context.Context) error {
//		for {
//			select {
//			case <-ctx.Done():
//				return nil
//			case item := <-items:
//				_ = item
//			}
//		}
//	}
//
// Returning nil for expected shutdown keeps graceful component stops distinct
// from task failures. A Task MAY still return context.Canceled,
// context.DeadlineExceeded, or context.Cause(ctx) when the context stop itself is
// the failure the component owner must observe.
//
// Package run does not recover panics raised by Task implementations. A panic is
// treated as a programming error owned by the task or by the component that
// supplied it. Owners that need panic-to-error conversion must wrap tasks
// explicitly outside the core Group contract.
type Task func(context.Context) error
