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
	return func(config *evaluatorConfig) error {
		policy = policy.Normalize()
		if err := policy.Validate(); err != nil {
			return err
		}

		config.executionPolicy = policy
		return nil
	}
}

// WithTargetExecutionPolicy configures the check execution policy for target.
//
// target MUST be concrete. TargetUnknown and invalid target values return an
// error classified as ErrInvalidTarget.
//
// Target-specific execution policies override the evaluator default execution
// policy. This is useful when readiness needs bounded parallel dependency probes
// while startup and liveness should remain sequential.
func WithTargetExecutionPolicy(target Target, policy ExecutionPolicy) EvaluatorOption {
	return func(config *evaluatorConfig) error {
		if !target.IsConcrete() {
			return InvalidTargetError{Target: target}
		}

		policy = policy.Normalize()
		if err := policy.Validate(); err != nil {
			return err
		}

		config.targetExecutionPolicies[target] = policy
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
func WithTargetSequentialChecks(target Target) EvaluatorOption {
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
func WithTargetParallelChecks(target Target, maxConcurrency int) EvaluatorOption {
	return WithTargetExecutionPolicy(target, ParallelExecutionPolicy(maxConcurrency))
}
