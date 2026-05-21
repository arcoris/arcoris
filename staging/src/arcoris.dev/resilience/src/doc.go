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


// Package resilience groups failure-control and execution-protection primitives
// for ARCORIS internals.
//
// Resilience packages compose lower-level building blocks into local protection
// behavior:
//
//   - retry owns bounded retry execution.
//   - deadline owns helpers for execution budgets derived from context
//     deadlines.
//   - retrybudget owns retry amplification limits.
//   - bulkhead owns bounded in-flight isolation.
//
// Some resilience packages also expose admission-compatible result surfaces for
// local primitives: bulkhead returns owned lease grants, retrybudget returns
// committed no-grant retry-spend decisions, and deadline returns no-side-effect
// start decisions. Those surfaces use admission contracts without moving global
// policy, catalog lookup, runtime chains, scheduling, health, metrics, logging,
// or tracing into the resilience root.
//
// chrono owns clocks, delays, jitter, and time primitives. runtime owns task,
// wait, and lifecycle mechanics. capacity owns local scalar capacity accounting.
// resilience composes those lower-level primitives into failure-control and
// execution-protection behavior.
//
// The root resilience layer must not import health or transport adapters.
// Health, metrics, logging, tracing, routing, scheduling, and admission policy
// are outside root resilience responsibility.
package resilience
