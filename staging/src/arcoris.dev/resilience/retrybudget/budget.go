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

package retrybudget

import "arcoris.dev/snapshot"

// OriginalRecorder records original, non-retry traffic that contributes to a
// retry budget.
//
// Implementations define what counts as original traffic. For a retry loop, this
// is usually the first operation attempt before any retry attempts are admitted.
type OriginalRecorder interface {
	// RecordOriginal records one original, non-retry attempt.
	RecordOriginal()
}

// RetryAdmitter atomically decides whether one retry attempt may be admitted.
//
// TryAdmitRetry is a check-and-spend operation. If it returns an allowed
// decision, the implementation has already accounted for the retry attempt. This
// avoids the race that would exist with separate CanRetry and RecordRetry calls.
type RetryAdmitter interface {
	// TryAdmitRetry decides whether one retry attempt may be admitted and returns
	// the resulting retry-budget decision.
	TryAdmitRetry() Decision
}

// Budget is the common contract for retry budget implementations.
//
// Budget embeds snapshot.Source[Snapshot] instead of defining a package-local
// Source interface. Revisioned snapshot publication belongs to arcoris.dev/snapshot;
// this package only defines the retry-budget domain value carried by that
// generic snapshot.
type Budget interface {
	OriginalRecorder
	RetryAdmitter
	snapshot.Source[Snapshot]
}
