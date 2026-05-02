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

import "errors"

var (
	// ErrNilRegistry identifies a nil Registry passed to NewEvaluator.
	//
	// Evaluator requires a registry owner. A nil registry would make target
	// evaluation ambiguous, so construction rejects it instead of producing an
	// evaluator that fails later.
	ErrNilRegistry = errors.New("health: nil registry")

	// ErrNilEvaluatorOption identifies a nil option passed to NewEvaluator.
	//
	// Options are executed during evaluator construction. A nil option is a
	// wiring error and is rejected explicitly.
	ErrNilEvaluatorOption = errors.New("health: nil evaluator option")

	// ErrNilClock identifies a nil clock passed to WithClock.
	//
	// Evaluator uses clock.PassiveClock for observation timestamps and elapsed
	// duration measurement. A nil clock would panic during evaluation.
	ErrNilClock = errors.New("health: nil clock")

	// ErrInvalidTimeout identifies a negative health check timeout.
	//
	// A zero timeout is valid and disables evaluator-enforced timeout. Negative
	// values do not describe a meaningful evaluation budget.
	ErrInvalidTimeout = errors.New("health: invalid timeout")
)
