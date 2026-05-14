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

package probe

import "reflect"

// NewRunner returns a Runner that periodically evaluates targets with evaluator.
//
// evaluator MUST be non-nil. Targets MUST be configured explicitly with
// WithTargets. NewRunner validates configuration and creates the private cache,
// but it does not start goroutines. Callers start probing by calling Run with an
// owner-controlled context.
func NewRunner(evaluator Evaluator, opts ...Option) (*Runner, error) {
	if nilEvaluator(evaluator) {
		return nil, ErrNilEvaluator
	}

	cfg := defaultConfig()
	if err := applyOptions(&cfg, opts...); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	targets := copyTargets(cfg.targets)

	return &Runner{
		evaluator:    evaluator,
		store:        newStore(targets),
		clock:        cfg.clock,
		targets:      targets,
		schedule:     cfg.schedule,
		staleAfter:   cfg.staleAfter,
		initialProbe: cfg.initialProbe,
	}, nil
}

func nilEvaluator(evaluator Evaluator) bool {
	if evaluator == nil {
		return true
	}

	value := reflect.ValueOf(evaluator)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
