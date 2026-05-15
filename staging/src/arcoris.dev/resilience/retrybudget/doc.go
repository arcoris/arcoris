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

// Package retrybudget defines contracts and shared value types for limiting
// retry amplification.
//
// A retry budget limits retry attempts relative to a budget signal such as
// original traffic, active work, or adaptive runtime capacity. This root package
// defines common recording, admission, decision, reason, kind, and snapshot
// contracts. Concrete accounting strategies live in implementation subpackages.
//
// Snapshots in this package are domain values. Revisioned publication and
// read-only access use arcoris.dev/snapshot. Implementations should expose
// snapshot.Source[Snapshot]. Mutable implementations may use
// snapshot.Publisher[Snapshot] or another snapshot source internally, depending
// on their ownership and performance model.
//
// The package does not execute retries, classify operation errors, compute
// delays, apply jitter, enforce deadlines, limit general request rate, provide
// circuit breaking, provide bulkheads, export metrics, integrate with health, or
// coordinate distributed/global budgets.
package retrybudget
