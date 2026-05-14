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

package eval

import (
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
)

// defaultCheckTimeout is the default timeout applied to checks by Evaluator.
//
// The default is intentionally finite so a stuck checker does not block the caller
// forever by default. Components with purely in-memory checks may explicitly
// disable the timeout with WithDefaultTimeout(0).
const defaultCheckTimeout = time.Second

// evaluatorConfig contains normalized Evaluator construction settings.
//
// The config is intentionally package-local. Public callers configure Evaluator
// through EvaluatorOption constructors, while NewEvaluator receives a complete
// normalized configuration.
//
// evaluatorConfig must stay small. If future health features need separate
// configuration domains, they should define their own option type rather than
// turning Evaluator into a global health configuration object.
type evaluatorConfig struct {
	// clock provides read-only runtime time.
	//
	// Evaluator depends on clock.PassiveClock because it only needs Now and Since
	// semantics. It does not own timers, tickers, sleeps, retry loops, probe
	// loops, or background scheduling behavior.
	clock clock.PassiveClock

	// defaultTimeout is the timeout applied to each check unless targetTimeouts
	// contains an override for the evaluated target.
	//
	// A zero value disables evaluator-enforced timeout. Negative values are
	// rejected by option constructors.
	defaultTimeout time.Duration

	// targetTimeouts stores per-target timeout overrides.
	//
	// The map is owned by evaluatorConfig during construction and is defensively
	// copied into Evaluator by NewEvaluator.
	targetTimeouts map[health.Target]time.Duration

	// executionPolicy is the default check execution policy used when
	// targetExecutionPolicies does not contain an override for the evaluated
	// target.
	//
	// The default is sequential. Parallel execution is always explicit and
	// bounded by ExecutionPolicy.MaxConcurrency.
	executionPolicy ExecutionPolicy

	// targetExecutionPolicies stores per-target execution policy overrides.
	//
	// The map is owned by evaluatorConfig during construction and is defensively
	// copied into Evaluator by NewEvaluator.
	targetExecutionPolicies map[health.Target]ExecutionPolicy
}

// defaultEvaluatorConfig returns the conservative Evaluator configuration.
//
// The default uses clock.RealClock, a one-second check timeout, and sequential
// check execution. The timeout is intentionally finite so a stuck checker does
// not block the caller forever by default. Components with purely in-memory
// checks may explicitly disable the timeout with WithDefaultTimeout(0).
func defaultEvaluatorConfig() evaluatorConfig {
	return evaluatorConfig{
		clock:                   clock.RealClock{},
		defaultTimeout:          defaultCheckTimeout,
		targetTimeouts:          make(map[health.Target]time.Duration),
		executionPolicy:         DefaultExecutionPolicy(),
		targetExecutionPolicies: make(map[health.Target]ExecutionPolicy),
	}
}
