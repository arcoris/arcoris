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

package eval

import (
	"time"

	"arcoris.dev/health"
)

// NewEvaluator returns an evaluator for resolver.
//
// resolver MUST be non-nil. The evaluator reads checks from the resolver at
// evaluation time. Resolver implementations define whether returned check sets
// are immutable snapshots, live views, or another deterministic lookup model.
//
// By default, every check receives a one-second timeout and checks are evaluated
// sequentially. Use WithDefaultTimeout to change or disable the default timeout,
// WithTargetTimeout to override the timeout for a specific target,
// WithExecutionPolicy or WithParallelChecks to change default execution, and
// target-specific execution options to override execution for a specific target.
func NewEvaluator(resolver health.CheckResolver, opts ...EvaluatorOption) (*Evaluator, error) {
	if resolver == nil {
		return nil, ErrNilResolver
	}

	cfg := defaultEvaluatorConfig()
	if err := applyEvaluatorOptions(&cfg, opts...); err != nil {
		return nil, err
	}

	return newEvaluatorFromConfig(resolver, cfg), nil
}

// newEvaluatorFromConfig copies normalized construction state into Evaluator.
func newEvaluatorFromConfig(resolver health.CheckResolver, cfg evaluatorConfig) *Evaluator {
	targetTimeouts := copyTargetTimeouts(cfg.targetTimeouts)
	targetExecutionPolicies := copyTargetExecutionPolicies(cfg.targetExecutionPolicies)

	return &Evaluator{
		resolver:                resolver,
		clock:                   cfg.clock,
		defaultTimeout:          cfg.defaultTimeout,
		targetTimeouts:          targetTimeouts,
		executionPolicy:         cfg.executionPolicy,
		targetExecutionPolicies: targetExecutionPolicies,
	}
}

// copyTargetTimeouts returns detached timeout overrides for Evaluator.
func copyTargetTimeouts(source map[health.Target]time.Duration) map[health.Target]time.Duration {
	copied := make(map[health.Target]time.Duration, len(source))
	for target, timeout := range source {
		copied[target] = timeout
	}

	return copied
}

// copyTargetExecutionPolicies returns detached execution-policy overrides.
func copyTargetExecutionPolicies(source map[health.Target]ExecutionPolicy) map[health.Target]ExecutionPolicy {
	copied := make(map[health.Target]ExecutionPolicy, len(source))
	for target, policy := range source {
		copied[target] = policy
	}

	return copied
}
