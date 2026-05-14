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

// Package resilience groups failure-control primitives for ARCORIS internals.
//
// The currently implemented package is retry. It owns bounded retry execution:
// attempts, retryability classification, retry-owned limits, context
// interruption, outcomes, observer events, and consumption of caller-provided
// delay.Schedule values.
//
// Deterministic delay formulas, randomized delay transforms, fake clocks, and
// timer abstractions belong to arcoris.dev/chrono. Runtime waiting mechanics
// belong to arcoris.dev/runtime. Health models and HTTP or gRPC adapters belong
// to arcoris.dev/health.
//
// Future failure-control siblings such as deadline, retrybudget,
// circuitbreaker, and bulkhead should live in this module when they are
// implemented. They are intentionally not added by this migration.
//
// resilience must not import health or transport adapter packages. retry depends
// only on the Go standard library plus arcoris.dev/chrono/clock and
// arcoris.dev/chrono/delay.
package resilience
