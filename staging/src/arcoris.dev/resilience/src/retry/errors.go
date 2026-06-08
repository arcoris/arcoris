// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import "errors"

var (
	// ErrNilContext reports a nil context at a public retry API boundary.
	ErrNilContext = errors.New("retry: nil context")

	// ErrNilOperation reports a nil Operation passed to Do or DoObserved.
	ErrNilOperation = errors.New("retry: nil operation")

	// ErrNilValueOperation reports a nil ValueOperation passed to DoValue,
	// DoValueObserved, or the private execution engine.
	ErrNilValueOperation = errors.New("retry: nil value operation")

	// ErrNilClock reports nil clock configuration.
	ErrNilClock = errors.New("retry: nil clock")

	// ErrNilDelaySchedule reports nil delay.Schedule configuration.
	ErrNilDelaySchedule = errors.New("retry: nil delay schedule")

	// ErrNilDelaySequence reports a delay.Schedule that returned a nil Sequence.
	ErrNilDelaySequence = errors.New("retry: delay schedule returned nil Sequence")

	// ErrNegativeDelay reports a delay.Sequence that returned a negative delay
	// while also reporting ok=true.
	ErrNegativeDelay = errors.New("retry: delay sequence returned negative delay")

	// ErrNilClassifier reports nil Classifier configuration.
	ErrNilClassifier = errors.New("retry: nil classifier")

	// ErrNilClassifierFunc reports a nil function adapted as a ClassifierFunc.
	ErrNilClassifierFunc = errors.New("retry: nil classifier function")

	// ErrZeroMaxAttempts reports an invalid zero max-attempt limit.
	ErrZeroMaxAttempts = errors.New("retry: zero max attempts")

	// ErrNegativeMaxElapsed reports an invalid negative elapsed-time limit.
	ErrNegativeMaxElapsed = errors.New("retry: negative max elapsed")

	// ErrNilObserver reports nil Observer configuration.
	ErrNilObserver = errors.New("retry: nil observer")

	// ErrNilObserverFunc reports a nil function adapted as an ObserverFunc.
	ErrNilObserverFunc = errors.New("retry: nil observer function")

	// ErrNilOption reports nil functional-option configuration.
	ErrNilOption = errors.New("retry: nil option")

	// ErrInvalidExhaustedOutcome reports invalid metadata passed to
	// NewExhaustedError.
	ErrInvalidExhaustedOutcome = errors.New("retry: invalid exhausted outcome")

	// ErrNonExhaustedOutcomeReason reports a valid but non-exhausted Outcome
	// passed to NewExhaustedError.
	ErrNonExhaustedOutcomeReason = errors.New("retry: non-exhausted outcome reason")
)
