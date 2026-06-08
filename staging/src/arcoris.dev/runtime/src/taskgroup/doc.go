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

// Package taskgroup provides context-scoped groups of named goroutine tasks.
//
// A Group owns one child context, a set of uniquely named Task functions, and a
// deterministic Wait result. It is intended for local component mechanics where
// an owner needs to start related goroutines, cancel siblings after failure, and
// collect named task errors without adding supervision, restart, retry, logging,
// metrics, tracing, worker-pool, queue, or lifecycle policy.
//
// A Group is constructor-created with New. Its zero value is invalid and all
// public methods panic when called on a nil or uninitialized Group. Public
// functions in this package also panic on nil contexts and nil functional
// options; callers that want a background context must pass context.Background
// explicitly.
//
// Task panics are not recovered. Panic recovery is owner policy and should be
// implemented by explicitly wrapping a Task before it is submitted.
//
// Example:
//
//	group := taskgroup.New(ctx)
//	group.Go("dispatcher", dispatcher.Run)
//	group.Go("controller", controller.Run)
//	err := group.Wait()
//
// AwaitContext is a convenience Task-compatible helper for groups that need a
// sentinel task bound only to a context. It is intentionally not named Wait so
// it cannot be confused with Group.Wait or the adjacent runtime/wait package.
//
// Production code in this package depends only on the Go standard library.
package taskgroup
