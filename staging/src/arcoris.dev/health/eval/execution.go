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

import "arcoris.dev/health"

// ExecutionMode identifies how Evaluator schedules checks during one target
// evaluation.
//
// ExecutionMode is evaluator-owned execution policy. It does not change checker
// contracts, registry semantics, result normalization, report aggregation,
// target policy, transport mapping, periodic probing, logging, metrics, tracing,
// retries, restart policy, admission policy, or scheduler behavior.
//
// The zero value is ExecutionSequential. This keeps Evaluator conservative and
// predictable unless a component owner explicitly opts into bounded parallel
// execution.
type ExecutionMode uint8

const (
	// ExecutionSequential evaluates checks one by one in health.Registry
	// registration order.
	//
	// Sequential execution is the default because it has the simplest load
	// profile, the smallest concurrency surface, and preserves the historical
	// Evaluator behavior.
	ExecutionSequential ExecutionMode = iota

	// ExecutionParallel evaluates checks concurrently with a bounded maximum
	// concurrency.
	//
	// Parallel execution is useful for independent I/O-bound checks such as
	// database, queue, cache, and storage probes. It MUST be explicitly
	// configured by the component owner because it may increase instantaneous
	// load on dependencies.
	ExecutionParallel
)

// String returns the stable diagnostic name of mode.
//
// The returned value is intended for diagnostics, tests, logs, and error
// messages. It is not a versioned wire format.
func (mode ExecutionMode) String() string {
	switch mode {
	case ExecutionSequential:
		return "sequential"
	case ExecutionParallel:
		return "parallel"
	default:
		return "invalid"
	}
}

// IsValid reports whether mode is one of the execution modes defined by this
// package.
func (mode ExecutionMode) IsValid() bool {
	switch mode {
	case ExecutionSequential,
		ExecutionParallel:
		return true
	default:
		return false
	}
}

// ExecutionPolicy configures how Evaluator executes all checks registered for
// one evaluated target.
//
// ExecutionPolicy only controls scheduling inside a single Evaluate call. It
// does not make Evaluate asynchronous: callers still wait for a complete health.Report.
// It also does not create timers, tickers, periodic probes, background loops,
// retries, caches, metrics, logs, traces, or transport adapters.
//
// The zero value normalizes to sequential execution. For parallel execution,
// MaxConcurrency MUST be positive.
type ExecutionPolicy struct {
	// Mode selects sequential or bounded parallel check execution.
	Mode ExecutionMode

	// MaxConcurrency limits the number of checks evaluated concurrently when Mode
	// is ExecutionParallel.
	//
	// Sequential execution normalizes MaxConcurrency to 1. Parallel execution
	// requires MaxConcurrency > 0. Evaluator may cap MaxConcurrency to the number
	// of checks during one evaluation.
	MaxConcurrency int
}

// DefaultExecutionPolicy returns the default evaluator execution policy.
//
// The default is sequential. Component owners may explicitly opt into bounded
// parallel execution globally or for a specific target through EvaluatorOption
// constructors.
func DefaultExecutionPolicy() ExecutionPolicy {
	return ExecutionPolicy{
		Mode:           ExecutionSequential,
		MaxConcurrency: 1,
	}
}

// ParallelExecutionPolicy returns a bounded parallel execution policy.
//
// The returned policy still needs validation because maxConcurrency <= 0 is not
// meaningful. Keeping construction and validation separate lets option
// constructors return typed validation errors with consistent classification.
func ParallelExecutionPolicy(maxConcurrency int) ExecutionPolicy {
	return ExecutionPolicy{
		Mode:           ExecutionParallel,
		MaxConcurrency: maxConcurrency,
	}
}

// Normalize returns a copy of policy with mode-specific default fields applied.
//
// Sequential execution always normalizes MaxConcurrency to 1 because no checks
// run concurrently. Parallel execution preserves MaxConcurrency so Validate can
// reject missing or invalid limits.
func (policy ExecutionPolicy) Normalize() ExecutionPolicy {
	switch policy.Mode {
	case ExecutionSequential:
		policy.MaxConcurrency = 1
	}

	return policy
}

// Validate reports whether policy is a supported evaluator execution policy.
func (policy ExecutionPolicy) Validate() error {
	policy = policy.Normalize()

	switch policy.Mode {
	case ExecutionSequential:
		return nil
	case ExecutionParallel:
		if policy.MaxConcurrency <= 0 {
			return InvalidExecutionPolicyError{
				Field:          "max_concurrency",
				Mode:           policy.Mode,
				MaxConcurrency: policy.MaxConcurrency,
			}
		}

		return nil
	default:
		return InvalidExecutionPolicyError{
			Field:          "mode",
			Mode:           policy.Mode,
			MaxConcurrency: policy.MaxConcurrency,
		}
	}
}

// WithExecutionPolicy configures the default check execution policy used by
// Evaluator.
//
// The default execution policy applies to every concrete target unless a
// target-specific execution policy is configured with WithTargetExecutionPolicy,
// WithTargetParallelChecks, or WithTargetSequentialChecks.
//
// The policy is normalized and validated at construction time. Invalid policies
// return an error classified as ErrInvalidExecutionPolicy.
func WithExecutionPolicy(policy ExecutionPolicy) EvaluatorOption {
	return func(cfg *evaluatorConfig) error {
		policy = policy.Normalize()
		if err := policy.Validate(); err != nil {
			return err
		}

		cfg.executionPolicy = policy
		return nil
	}
}

// WithTargetExecutionPolicy configures the check execution policy for target.
//
// target MUST be concrete. health.TargetUnknown and invalid target values return
// an error classified as health.ErrInvalidTarget.
//
// Target-specific execution policies override the evaluator default execution
// policy. This is useful when readiness needs bounded parallel dependency probes
// while startup and liveness should remain sequential.
func WithTargetExecutionPolicy(target health.Target, policy ExecutionPolicy) EvaluatorOption {
	return func(cfg *evaluatorConfig) error {
		if !target.IsConcrete() {
			return health.InvalidTargetError{Target: target}
		}

		policy = policy.Normalize()
		if err := policy.Validate(); err != nil {
			return err
		}

		cfg.targetExecutionPolicies[target] = policy
		return nil
	}
}

// WithSequentialChecks configures Evaluator to execute checks sequentially by
// default.
//
// Target-specific execution policies still override this default.
func WithSequentialChecks() EvaluatorOption {
	return WithExecutionPolicy(DefaultExecutionPolicy())
}

// WithTargetSequentialChecks configures Evaluator to execute checks for target
// sequentially.
//
// This option is useful when a global parallel policy is configured but a
// specific target should retain the conservative sequential execution model.
func WithTargetSequentialChecks(target health.Target) EvaluatorOption {
	return WithTargetExecutionPolicy(target, DefaultExecutionPolicy())
}

// WithParallelChecks configures Evaluator to execute checks with bounded
// parallelism by default.
//
// maxConcurrency MUST be positive. Target-specific execution policies still
// override this default.
func WithParallelChecks(maxConcurrency int) EvaluatorOption {
	return WithExecutionPolicy(ParallelExecutionPolicy(maxConcurrency))
}

// WithTargetParallelChecks configures Evaluator to execute checks for target with
// bounded parallelism.
//
// target MUST be concrete. maxConcurrency MUST be positive. The configured
// policy applies only to target and overrides the evaluator default execution
// policy.
func WithTargetParallelChecks(target health.Target, maxConcurrency int) EvaluatorOption {
	return WithTargetExecutionPolicy(target, ParallelExecutionPolicy(maxConcurrency))
}
