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

// applyEvaluatorOptions applies options to normalized configuration in order.
//
// Later options win for single-value domains. Target timeout options replace the
// previous timeout for the same target. Target execution options replace the
// previous execution policy for the same target. Nil options are rejected with
// ErrNilEvaluatorOption so invalid conditional option composition is visible at
// the construction boundary.
func applyEvaluatorOptions(cfg *evaluatorConfig, opts ...EvaluatorOption) error {
	for _, opt := range opts {
		if opt == nil {
			return ErrNilEvaluatorOption
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}

	return nil
}
