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
	"context"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
)

// Evaluator executes registered health checks and returns target reports.
//
// Evaluator is transport-neutral. It does not expose HTTP handlers, map gRPC
// serving states, log diagnostics, emit metrics, perform retries, run periodic
// probes, or decide restart, admission, routing, or scheduling behavior. It only
// owns the synchronous evaluation boundary for checks registered in health.Registry.
//
// Evaluation is deterministic with respect to registry order. The default
// execution policy is sequential. Component owners may configure bounded
// parallel execution globally or per target. Parallel execution preserves
// health.Report.Checks order by registry registration order even when checks finish out
// of order.
//
// Evaluator applies a cooperative context to every check and, when a timeout is
// configured, also enforces a caller-visible result boundary. A checker that
// ignores its context may continue running after the evaluator has returned a
// timeout result. health.Checker implementations SHOULD observe ctx whenever they can
// block, perform I/O, wait on another goroutine, or acquire external resources.
//
// Evaluator recovers checker panics and converts them into unhealthy results
// with health.ReasonPanic. Panic details are preserved only in health.Result.Cause and MUST
// NOT be exposed by public adapters by default.
type Evaluator struct {
	registry *health.Registry

	clock          clock.PassiveClock
	defaultTimeout time.Duration
	targetTimeouts map[health.Target]time.Duration

	executionPolicy         ExecutionPolicy
	targetExecutionPolicies map[health.Target]ExecutionPolicy
}

// NewEvaluator returns an evaluator for registry.
//
// registry MUST be non-nil. The evaluator reads checks from the registry at
// evaluation time, so later registry mutations are visible to later evaluations.
// Component owners SHOULD normally finish registration during setup and treat
// the registry as effectively immutable after startup.
//
// By default, every check receives a one-second timeout and checks are evaluated
// sequentially. Use WithDefaultTimeout to change or disable the default timeout,
// WithTargetTimeout to override the timeout for a specific target,
// WithExecutionPolicy or WithParallelChecks to change default execution, and
// target-specific execution options to override execution for a specific target.
func NewEvaluator(r *health.Registry, opts ...EvaluatorOption) (*Evaluator, error) {
	if r == nil {
		return nil, ErrNilRegistry
	}

	cfg := defaultEvaluatorConfig()
	if err := applyEvaluatorOptions(&cfg, opts...); err != nil {
		return nil, err
	}

	targetTimeouts := make(map[health.Target]time.Duration, len(cfg.targetTimeouts))
	for target, timeout := range cfg.targetTimeouts {
		targetTimeouts[target] = timeout
	}

	targetExecutionPolicies := make(map[health.Target]ExecutionPolicy, len(cfg.targetExecutionPolicies))
	for target, policy := range cfg.targetExecutionPolicies {
		targetExecutionPolicies[target] = policy
	}

	return &Evaluator{
		registry:                r,
		clock:                   cfg.clock,
		defaultTimeout:          cfg.defaultTimeout,
		targetTimeouts:          targetTimeouts,
		executionPolicy:         cfg.executionPolicy,
		targetExecutionPolicies: targetExecutionPolicies,
	}, nil
}

// Evaluate runs all checks registered for target and returns an aggregated
// report.
//
// target MUST be concrete. Invalid or non-concrete targets return a health.StatusUnknown
// report and an error classified as health.ErrInvalidTarget.
//
// A nil ctx is treated as context.Background. This mirrors defensive boundaries
// in other adjacent ARCORIS packages and avoids panics in diagnostics and tests.
//
// If target has no registered checks, Evaluate returns a health.StatusUnknown report.
// Absence of checks is not treated as healthy because health requires an
// affirmative observation.
//
// Evaluate is synchronous regardless of execution policy. Parallel execution
// affects only how checks are scheduled inside this call; the caller still
// receives one complete health.Report after all scheduled checks have produced
// caller-visible Results.
func (e *Evaluator) Evaluate(ctx context.Context, target health.Target) (health.Report, error) {
	started := e.clock.Now()

	if !target.IsConcrete() {
		return health.Report{
			Target:   target,
			Status:   health.StatusUnknown,
			Observed: started,
		}, health.InvalidTargetError{Target: target}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	checks := e.registry.Checks(target)
	if len(checks) == 0 {
		return health.Report{
			Target:   target,
			Status:   health.StatusUnknown,
			Observed: started,
		}, nil
	}

	timeout := e.timeoutFor(target)
	executionPolicy := e.executionPolicyFor(target)
	results := e.evaluateChecks(ctx, checks, timeout, executionPolicy)
	status := aggregateStatus(results)

	finished := e.clock.Now()

	return health.Report{
		Target:   target,
		Status:   status,
		Observed: finished,
		Duration: nonNegativeDuration(e.clock.Since(started)),
		Checks:   results,
	}, nil
}

// timeoutFor returns the effective check timeout for target.
func (e *Evaluator) timeoutFor(target health.Target) time.Duration {
	if timeout, ok := e.targetTimeouts[target]; ok {
		return timeout
	}

	return e.defaultTimeout
}

// executionPolicyFor returns the effective check execution policy for target.
//
// health.Target-specific execution policy overrides the evaluator default. The returned
// policy is normalized at construction time.
func (e *Evaluator) executionPolicyFor(target health.Target) ExecutionPolicy {
	if policy, ok := e.targetExecutionPolicies[target]; ok {
		return policy
	}

	return e.executionPolicy
}
