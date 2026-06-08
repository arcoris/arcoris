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
// A Limiter must be created with New. Nil or uninitialized Limiter receivers
// panic with package-owned receiver-validation errors.
//
// Fixed windows are local and observation-aligned. The first window starts when
// the limiter is created, and subsequent windows start when a write path observes
// that the previous window has ended. Windows are not aligned to wall-clock
// boundaries such as minute, hour, or day boundaries.
//
// Snapshot reads do not rotate windows. A quiet limiter may keep publishing the
// last observed window until RecordOriginal or TryAdmitRetry observes time
// advancement.
//
// Default retry capacity uses exact ratio configuration:
//
//	allowed = minRetries + floor(originalAttempts * ratio.Numerator / ratio.Denominator)
//
// The result saturates instead of wrapping. The ratio is intentionally bounded to
// the conservative range [0, 1]. Ratio math is integer-only and remains exact
// for the full uint64 original-attempt range.
//
// Minimum retry allowance is available at the start of each window, even before
// RecordOriginal observes traffic. Set minRetries to zero for strict
// traffic-proportional behavior.
//
// Limiter is local to one process. It does not coordinate distributed budgets,
// smooth window boundaries, execute retries, classify errors, compute delays,
// enforce deadlines, provide a circuit breaker, provide a bulkhead, export
// metrics, or integrate with health. Those responsibilities belong to other
// resilience components or adapter packages.
package fixedwindow
