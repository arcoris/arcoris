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

// Package fixedwindow provides a local fixed-window retry budget implementation.
//
// The limiter admits retry attempts according to a simple traffic-ratio budget
// within one local time window:
//
//	allowed = minRetries + floor(originalAttempts * ratio)
//
// Original attempts are recorded with RecordOriginal. Retry attempts are admitted
// with TryAdmitRetry, which is an atomic check-and-spend operation. A successful
// admission records the retry attempt before returning to the caller.
//
// Limiter is local to one process. It does not coordinate distributed budgets,
// smooth window boundaries, execute retries, classify errors, compute delays,
// enforce deadlines, provide a circuit breaker, provide a bulkhead, export
// metrics, or integrate with health. Those responsibilities belong to other
// resilience components or adapter packages.
package fixedwindow
