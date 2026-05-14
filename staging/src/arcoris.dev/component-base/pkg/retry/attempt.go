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

package retry

import "time"

// Attempt describes one operation call performed by a retry execution.
//
// Attempt is immutable retry metadata. It records the one-based attempt number
// and the retry clock time at which that attempt started. It does not execute
// work, classify errors, compute retry delays, wait on timers, observe
// contexts, emit events, or decide whether another attempt is allowed.
//
// Attempt values are copyable. They must not contain locks, channels, mutable
// delay sequences, observers, operations, or references back to the retry
// execution that produced them.
//
// Attempt numbering is one-based:
//
//   - attempt 1 is the initial operation call;
//   - attempt 2 and later are retry operation calls;
//   - attempt 0 is invalid and represents no operation call.
//
// StartedAt is assigned by the retry execution immediately before the operation
// is called. A zero StartedAt means the value does not describe a started
// attempt and is therefore not valid attempt metadata.
//
// Attempt intentionally does not store the operation error or the retry delay.
// Errors and delays belong to retry events and outcomes, not to the identity of
// the operation call itself.
type Attempt struct {
	// Number is the one-based operation call number within one retry execution.
	//
	// Number 1 is the initial operation call. Numbers greater than 1 are retry
	// operation calls. Number 0 is invalid and must not be used for a started
	// attempt.
	Number uint

	// StartedAt is the retry clock time at which this attempt started.
	//
	// The value is assigned before the operation is called. It is intended for
	// diagnostics, observer events, elapsed-time accounting, and tests. It is not
	// a distributed ordering source and must not be interpreted as a cluster-wide
	// logical clock.
	StartedAt time.Time
}

// IsZero reports whether a is the zero Attempt value.
//
// The zero Attempt is not a valid started attempt. It is useful as an omitted
// value in events or outcomes that can occur before any operation call is made,
// such as retry-owned context interruption before the first attempt.
func (a Attempt) IsZero() bool {
	return a.Number == 0 && a.StartedAt.IsZero()
}

// IsValid reports whether a describes a started operation attempt.
//
// A valid Attempt has a non-zero one-based Number and a non-zero StartedAt time.
// IsValid does not check whether StartedAt is in the past, present, or future
// relative to any particular clock. Retry executions own clock selection and
// timestamp assignment.
func (a Attempt) IsValid() bool {
	return a.Number != 0 && !a.StartedAt.IsZero()
}

// IsFirst reports whether a describes the initial operation call.
//
// IsFirst returns false for invalid attempts, including values with Number 1 but
// a zero StartedAt timestamp. Use IsValid when callers need to distinguish
// malformed attempt metadata from a valid first attempt.
func (a Attempt) IsFirst() bool {
	return a.IsValid() && a.Number == 1
}

// IsRetry reports whether a describes a retry operation call.
//
// Retry attempts are operation calls after the initial attempt. IsRetry returns
// false for invalid attempts, including values with Number greater than 1 but a
// zero StartedAt timestamp.
func (a Attempt) IsRetry() bool {
	return a.IsValid() && a.Number > 1
}
