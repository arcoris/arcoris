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

// Package chrono groups time abstraction and duration-sequence construction
// modules for ARCORIS internals.
//
// The module owns clock, delay, and jitter. Package clock provides real and fake
// time sources, timers, and tickers. Package delay provides deterministic delay
// schedules and wrappers. Package jitter provides non-cryptographic randomized
// delay schedules and randomized wrappers over delay.Schedule.
//
// chrono intentionally does not own runtime waiting mechanics. Blocking waits,
// condition loops, wait-owned context errors, and owner-controlled real-runtime
// timers belong to arcoris.dev/runtime/wait. Retry execution and other
// failure-control primitives belong to arcoris.dev/resilience. Health checks,
// probes, and transport adapters belong to arcoris.dev/health.
//
// Production code in clock and delay depends only on the Go standard library.
// jitter may depend on chrono/delay. The module must not depend on runtime,
// resilience, health, or transport adapter packages.
package chrono
