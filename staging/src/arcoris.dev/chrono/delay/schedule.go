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

// Schedule describes a reusable recipe for producing delay sequences.
//
// Schedule is the top-level contract of this package. It does not represent one
// running retry loop, polling loop, reconnect loop, or cooldown sequence.
// Instead, it describes how such a loop should create its own independent stream
// of delay values.
//
// The separation between Schedule and Sequence is intentional:
//
//   - Schedule is a reusable, policy-neutral delay recipe;
//   - Sequence is a per-owner stream of concrete delay values;
//   - retry packages own attempts, retryability decisions, errors, and context;
//   - clock packages own timers, sleeps, tickers, and fake-time behavior;
//   - wait packages own low-level wait mechanics and condition loops.
//
// This package intentionally keeps Schedule limited to delay generation. A
// Schedule MUST NOT sleep, create timers, observe contexts, execute operations,
// classify errors, retry failed work, rate limit callers, schedule queue items,
// log, trace, export metrics, or make domain-specific decisions.
//
// Schedule values returned by package constructors SHOULD be immutable after
// construction and SHOULD be safe for concurrent calls to NewSequence. This
// allows a retry policy, controller policy, or reconnect policy to store one
// Schedule and reuse it for many independent executions.
//
// Each call to NewSequence SHOULD return an independent Sequence. Sharing one
// stateful Sequence across unrelated loops is usually incorrect because
// sequences may track attempt index, previous delay, finite exhaustion state, or
// deterministic random state.
//
// A Schedule may produce finite or infinite sequences. Finite sequence
// exhaustion is represented by Sequence.Next returning ok=false. Exhaustion is
// not an error at the delay layer. The owner decides whether exhaustion means
// retry exhaustion, polling exhaustion, fallback, shutdown, or another
// higher-level outcome.
type Schedule interface {
	// NewSequence creates a delay sequence for one runtime owner.
	//
	// The returned Sequence MUST be non-nil. Implementations SHOULD return a
	// fresh independent sequence for each call unless the sequence is stateless
	// or is explicitly documented as safe to share.
	NewSequence() Sequence
}
