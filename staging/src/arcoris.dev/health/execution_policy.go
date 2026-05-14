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

// ExecutionPolicy configures how Evaluator executes all checks registered for
// one evaluated target.
//
// ExecutionPolicy only controls scheduling inside a single Evaluate call. It
// does not make Evaluate asynchronous: callers still wait for a complete Report.
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
