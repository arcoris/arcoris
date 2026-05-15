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

// Package noop provides an unlimited retry-budget implementation.
//
// A noop Budget admits every retry attempt and does not record original or retry
// traffic. It is useful when retry-budget enforcement is intentionally disabled
// but callers still want to depend on the retrybudget.Budget contract without
// nil checks or special cases.
//
// The implementation is stateless, concurrency-safe, zero-value usable, and
// allocation-free on the read path. Snapshot returns a stable revisioned snapshot
// through arcoris.dev/snapshot. No snapshot.Publisher is needed because the
// published state is immutable and never changes.
package noop
