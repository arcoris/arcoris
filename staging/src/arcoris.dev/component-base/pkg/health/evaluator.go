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
	"context"
	"time"

	"arcoris.dev/component-base/pkg/clock"
)

// Evaluator executes registered health checks and returns target reports.
//
// Evaluator is transport-neutral. It does not expose HTTP handlers, map gRPC
// serving states, log diagnostics, emit metrics, perform retries, run periodic
// probes, or decide restart, admission, routing, or scheduling behavior. It only
// owns the synchronous evaluation boundary for checks registered in Registry.
//
// Evaluation is deterministic with respect to registry order. Checks are
// evaluated sequentially in the order they were registered for the target.
//
// Evaluator applies a cooperative context to every check and, when a timeout is
// configured, also enforces a caller-visible result boundary. A checker that
// ignores its context may continue running after the evaluator has returned a
// timeout result. Checker implementations SHOULD observe ctx whenever they can
// block, perform I/O, wait on another goroutine, or acquire external resources.
//
// Evaluator recovers checker panics and converts them into unhealthy results
// with ReasonPanic. Panic details are preserved only in Result.Cause and MUST
// NOT be exposed by public adapters by default.
type Evaluator struct {
	registry *Registry

	clock          clock.PassiveClock
	defaultTimeout time.Duration
	targetTimeouts map[Target]time.Duration
}

// NewEvaluator returns an evaluator for registry.
//
// registry MUST be non-nil. The evaluator reads checks from the registry at
// evaluation time, so later registry mutations are visible to later evaluations.
// Component owners SHOULD normally finish registration during setup and treat
// the registry as effectively immutable after startup.
//
// By default, every check receives a one-second timeout. Use WithDefaultTimeout
// to change or disable the default timeout, and WithTargetTimeout to override the
// timeout for a specific target.
func NewEvaluator(registry *Registry, opts ...EvaluatorOption) (*Evaluator, error) {
	if registry == nil {
		return nil, ErrNilRegistry
	}

	cfg := defaultEvaluatorConfig()
	if err := applyEvaluatorOptions(&cfg, opts...); err != nil {
		return nil, err
	}

	targetTimeouts := make(map[Target]time.Duration, len(cfg.targetTimeouts))
	for target, timeout := range cfg.targetTimeouts {
		targetTimeouts[target] = timeout
	}

	return &Evaluator{
		registry:       registry,
		clock:          cfg.clock,
		defaultTimeout: cfg.defaultTimeout,
		targetTimeouts: targetTimeouts,
	}, nil
}

// Evaluate runs all checks registered for target and returns an aggregated
// report.
//
// target MUST be concrete. Invalid or non-concrete targets return a StatusUnknown
// report and an error classified as ErrInvalidTarget.
//
// A nil ctx is treated as context.Background. This mirrors defensive boundaries
// in other component-base packages and avoids panics in diagnostics and tests.
//
// If target has no registered checks, Evaluate returns a StatusUnknown report.
// Absence of checks is not treated as healthy because health requires an
// affirmative observation.
func (e *Evaluator) Evaluate(ctx context.Context, target Target) (Report, error) {
	started := e.clock.Now()

	if !target.IsConcrete() {
		return Report{
			Target:   target,
			Status:   StatusUnknown,
			Observed: started,
		}, InvalidTargetError{Target: target}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	checks := e.registry.Checks(target)
	if len(checks) == 0 {
		return Report{
			Target:   target,
			Status:   StatusUnknown,
			Observed: started,
		}, nil
	}

	results := make([]Result, 0, len(checks))
	status := StatusHealthy
	timeout := e.timeoutFor(target)

	for _, check := range checks {
		result := e.evaluateCheck(ctx, check, timeout)

		results = append(results, result)
		if result.Status.MoreSevereThan(status) {
			status = result.Status
		}
	}

	finished := e.clock.Now()

	return Report{
		Target:   target,
		Status:   status,
		Observed: finished,
		Duration: nonNegativeDuration(finished.Sub(started)),
		Checks:   results,
	}, nil
}

// timeoutFor returns the effective check timeout for target.
func (e *Evaluator) timeoutFor(target Target) time.Duration {
	if timeout, ok := e.targetTimeouts[target]; ok {
		return timeout
	}

	return e.defaultTimeout
}
