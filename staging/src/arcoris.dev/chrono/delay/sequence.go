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

package delay

import "time"

// Sequence is a single-owner stream of delay values.
//
// Sequence is the mutable counterpart to Schedule. Schedule values describe how
// delays should be generated. Sequence values produce the concrete delay stream
// for one runtime owner, such as one retry loop, polling loop, reconnect loop, or
// controller cooldown path.
//
// A Sequence typically records iteration state such as:
//
//   - the next index in an explicit delay list;
//   - the current step in a linear or exponential schedule;
//   - the previous delay in a decorrelated-jitter schedule;
//   - the number of remaining values in a finite schedule;
//   - per-sequence pseudo-random state for deterministic jitter tests.
//
// Sequence values intentionally do not own runtime waiting. They MUST NOT sleep,
// create timers, observe contexts, execute operations, classify errors, retry
// failed work, log, trace, export metrics, rate limit callers, schedule queue
// items, or make scheduler, admission, lifecycle, or domain decisions. They only
// return duration values.
//
// Sequence values are single-owner by default. Implementations are not required
// to be safe for concurrent calls to Next unless a concrete implementation
// explicitly documents stronger guarantees. Avoiding mandatory synchronization
// keeps hot retry and polling paths cheap and leaves ownership explicit.
//
// A zero delay is valid and means the owner may continue immediately. A negative
// delay is invalid and MUST NOT be produced by implementations in this package.
// If duration arithmetic would overflow, implementations SHOULD saturate or cap
// according to their documented algorithm instead of returning wrapped negative
// values.
//
// A Sequence may be finite or infinite. Finite sequence exhaustion is reported by
// Next returning ok=false. Exhaustion is not an error at the delay layer. The
// owner decides whether exhaustion means retry exhaustion, polling exhaustion,
// fallback, shutdown, or another higher-level outcome.
type Sequence interface {
	// Next returns the next delay in the sequence.
	//
	// The returned delay is meaningful only when ok is true.
	//
	// A result of delay=0, ok=true means the owner may continue immediately. A
	// result of ok=false means the sequence is exhausted and no further delay is
	// available.
	//
	// Implementations in this package must return non-negative delays when ok is
	// true. A negative delay with ok=true violates the Sequence contract.
	Next() (delay time.Duration, ok bool)
}
