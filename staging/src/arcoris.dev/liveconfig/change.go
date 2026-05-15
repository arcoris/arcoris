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

package liveconfig

import "arcoris.dev/snapshot"

// Change describes the result of applying a candidate configuration.
//
// Previous is the snapshot that was current before Apply evaluated the
// candidate. Current is the snapshot that remains current after Apply returns.
//
// Changed and Reason are intentionally separate dimensions. Changed answers the
// fast state question: "did this Apply publish a new source revision?" Reason
// answers the diagnostic question: "why did this Apply publish, skip
// publication, or reject the candidate?" Callers that only care about cache
// invalidation or cheap reload decisions can keep using Changed. Callers that
// need logs, counters, or operator diagnostics should use Reason.
//
// Reason is classification, not detail. When Reason is a failure value, the
// error returned by Apply carries the detailed normalization or validation
// failure and LastError records it. Change deliberately does not contain an
// error field so there is only one detailed error channel for Apply.
//
// The holder maintains these invariants:
//
//   - ChangeReasonPublished: Changed is true, Apply returns nil error, and
//     Current.Revision differs from Previous.Revision.
//   - ChangeReasonEqual: Changed is false, Apply returns nil error, and
//     Current.Revision equals Previous.Revision.
//   - ChangeReasonNormalizeFailed: Changed is false, Apply returns a non-nil
//     error, and Current.Revision equals Previous.Revision.
//   - ChangeReasonValidateFailed: Changed is false, Apply returns a non-nil
//     error, and Current.Revision equals Previous.Revision.
type Change[T any] struct {
	// Previous is the snapshot visible before the candidate was evaluated.
	Previous snapshot.Snapshot[T]

	// Current is the snapshot visible after the candidate was evaluated.
	//
	// For rejected candidates and accepted equal candidates, Current is the same
	// source revision as Previous. For published candidates, Current is the newly
	// published snapshot.
	Current snapshot.Snapshot[T]

	// Changed reports whether Apply published a new source revision.
	//
	// Changed is true only when Reason is ChangeReasonPublished.
	Changed bool

	// Reason classifies why Apply published, skipped publication, or rejected the
	// candidate.
	//
	// Reason is suitable for stable diagnostics. It should not be used as a
	// substitute for the detailed error returned by Apply when Rejected is true.
	Reason ChangeReason
}

// IsChanged reports whether Apply published a new source revision.
//
// It is equivalent to reading Changed directly and is retained for callers that
// prefer method-style predicates.
func (c Change[T]) IsChanged() bool {
	return c.Changed
}

// IsNoop reports whether Apply left the current source revision unchanged.
//
// This includes both accepted equal candidates and rejected candidates. Use
// Reason, Accepted, or Rejected when diagnostics need to distinguish them.
func (c Change[T]) IsNoop() bool {
	return !c.Changed
}

// Accepted reports whether Apply accepted the candidate.
//
// Published and equal candidates are both accepted. Equal candidates are valid
// no-ops, not rejected updates.
func (c Change[T]) Accepted() bool {
	return c.Reason.Accepted()
}

// Rejected reports whether Apply rejected the candidate.
//
// Rejected changes preserve the previous last-good config and correspond to the
// non-nil error returned by Apply.
func (c Change[T]) Rejected() bool {
	return c.Reason.Rejected()
}
