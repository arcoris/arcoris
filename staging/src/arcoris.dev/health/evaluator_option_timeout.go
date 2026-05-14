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

import "time"

// WithDefaultTimeout configures the default timeout applied to every check.
//
// A zero timeout disables evaluator-enforced timeout. This is useful for purely
// in-memory checks or for tests that want to avoid goroutine-based timeout
// boundaries.
//
// A negative timeout returns ErrInvalidTimeout.
//
// Target-specific timeouts configured with WithTargetTimeout override the
// default timeout for that target.
func WithDefaultTimeout(timeout time.Duration) EvaluatorOption {
	return func(cfg *evaluatorConfig) error {
		if timeout < 0 {
			return ErrInvalidTimeout
		}

		cfg.defaultTimeout = timeout
		return nil
	}
}

// WithTargetTimeout configures the timeout applied to checks for target.
//
// target MUST be concrete. TargetUnknown and invalid target values return an
// error classified as ErrInvalidTarget.
//
// A zero timeout disables evaluator-enforced timeout for the target. A negative
// timeout returns ErrInvalidTimeout.
//
// Target-specific timeouts are useful when startup, liveness, and readiness have
// different evaluation budgets. For example, liveness checks should normally be
// cheap and short, while startup checks may need a wider budget during bootstrap.
func WithTargetTimeout(target Target, timeout time.Duration) EvaluatorOption {
	return func(cfg *evaluatorConfig) error {
		if !target.IsConcrete() {
			return InvalidTargetError{Target: target}
		}
		if timeout < 0 {
			return ErrInvalidTimeout
		}

		cfg.targetTimeouts[target] = timeout
		return nil
	}
}
