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

package health

import (
	"time"

	"arcoris.dev/chrono/clock"
)

const defaultCheckTimeout = time.Second

// EvaluatorOption configures an Evaluator at construction time.
//
// Options are applied by NewEvaluator to a private evaluatorConfig before the
// Evaluator is created. They do not mutate an already constructed Evaluator and
// are not retained as option values after configuration normalization.
//
// Evaluator options must remain limited to evaluator-owned mechanics:
//
//   - time source for observation timestamps and duration measurement;
//   - default check timeout;
//   - target-specific check timeouts;
//   - default check execution policy;
//   - target-specific check execution policies.
//
// Evaluator options must not configure HTTP paths, gRPC serving states, metrics,
// tracing, logging, restart policy, admission policy, scheduler policy,
// dependency-specific checks, or periodic probe execution.
type EvaluatorOption func(*evaluatorConfig) error

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
	targetTimeouts map[Target]time.Duration

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
	targetExecutionPolicies map[Target]ExecutionPolicy
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
		targetTimeouts:          make(map[Target]time.Duration),
		executionPolicy:         DefaultExecutionPolicy(),
		targetExecutionPolicies: make(map[Target]ExecutionPolicy),
	}
}

// applyEvaluatorOptions applies options to config in order.
//
// Later options win for single-value domains. Target timeout options replace the
// previous timeout for the same target. Target execution options replace the
// previous execution policy for the same target. Nil options are rejected with
// ErrNilEvaluatorOption so invalid conditional option composition is visible at
// the construction boundary.
func applyEvaluatorOptions(config *evaluatorConfig, options ...EvaluatorOption) error {
	for _, option := range options {
		if option == nil {
			return ErrNilEvaluatorOption
		}
		if err := option(config); err != nil {
			return err
		}
	}

	return nil
}
