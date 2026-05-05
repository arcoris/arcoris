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

package healthprobe

import "arcoris.dev/component-base/pkg/health"

// NewRunner returns a Runner that periodically evaluates targets with evaluator.
//
// evaluator MUST be non-nil. Targets MUST be configured explicitly with
// WithTargets. NewRunner validates configuration and creates the private cache,
// but it does not start goroutines. Callers start probing by calling Run with an
// owner-controlled context.
func NewRunner(evaluator *health.Evaluator, options ...Option) (*Runner, error) {
	if evaluator == nil {
		return nil, ErrNilEvaluator
	}

	config := defaultConfig()
	if err := applyOptions(&config, options...); err != nil {
		return nil, err
	}
	if err := config.validate(); err != nil {
		return nil, err
	}

	targets := copyTargets(config.targets)

	return &Runner{
		evaluator:    evaluator,
		store:        newStore(targets),
		clock:        config.clock,
		targets:      targets,
		interval:     config.interval,
		staleAfter:   config.staleAfter,
		initialProbe: config.initialProbe,
	}, nil
}
