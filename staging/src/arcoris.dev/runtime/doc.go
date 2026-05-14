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

// Package runtime groups process-local runtime coordination primitives for
// ARCORIS internals.
//
// The module owns context-first task groups, real-runtime waiting mechanics,
// process signal integration, and component lifecycle state coordination through
// the run, wait, signals, and lifecycle packages.
//
// wait belongs here because it owns blocking runtime mechanics: cancellable
// delays, owner-controlled timers, condition loops, positive real-runtime jitter,
// and wait-owned context error classification. It is intentionally separate from
// arcoris.dev/chrono, which owns time sources and duration sequence
// construction.
//
// runtime does not own retry execution, retry budgets, circuit breakers, health
// models, health transports, delay formulas, jitter algorithms, schedulers,
// queues, worker pools, or observability backends. Failure-control primitives
// belong to arcoris.dev/resilience. Health checks and adapters belong to
// arcoris.dev/health.
//
// Production packages in this module should remain standard-library-only unless
// a concrete runtime package has a narrow reason to depend on another focused
// ARCORIS module. The module must not depend on resilience or health.
package runtime
