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

import "context"

// Observer receives retry execution events.
//
// Observer is a notification boundary. Retry execution calls configured
// observers when attempt, delay, and terminal stop events occur. Observers may
// record diagnostics, logs, metrics, traces, tests, or other caller-owned
// side-channel data.
//
// Observer must not be used as retry policy. It does not decide whether an error
// is retryable, does not choose retry delays, does not change attempt limits,
// does not mutate retry configuration, and does not control operation execution.
//
// ObserveRetry does not return an error. Observer failures cannot change retry
// execution. Implementations that write to logging, metrics, tracing, or other
// external systems must handle their own internal failures according to their
// adapter policy.
//
// Retry execution calls observers synchronously. Implementations SHOULD be fast
// and SHOULD avoid blocking work. Slow observers delay the retry loop that calls
// them.
//
// The retry package does not recover panics raised by observers. Panic recovery,
// if required, belongs to the observer implementation, runtime supervisor, or an
// explicit wrapper outside this package.
//
// Observer implementations must be safe for the way they are configured. If the
// same observer instance is shared by concurrent retry executions, the
// implementation must provide its own synchronization.
//
// ObserveRetry receives the retry-owned context passed to the retry execution.
// Observers MAY use the context for correlation or cancellation-aware side work,
// but they MUST NOT assume they can change the retry decision by modifying or
// observing the context.
//
// Event validation is owned by retry execution and tests. Observer
// implementations MAY defensively check event.IsValid, but the Observer contract
// itself does not require every implementation to reject malformed events.
type Observer interface {
	// ObserveRetry observes one retry event.
	//
	// The event value is immutable metadata for the observer. Implementations must
	// not assume that retaining ctx after ObserveRetry returns is safe unless they
	// own the resulting lifetime and cancellation behavior explicitly.
	ObserveRetry(ctx context.Context, event Event)
}
