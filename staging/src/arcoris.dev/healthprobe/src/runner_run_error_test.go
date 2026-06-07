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

package probe

import (
	"context"
	"errors"
	"testing"
)

func TestRunnerRunRejectsNilSequence(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(
		t,
		newTestClock(),
		WithSchedule(nilSequenceSchedule{}),
		WithInitialProbe(false),
	)

	err := runner.Run(context.Background())

	if !errors.Is(err, ErrNilSequence) {
		t.Fatalf("Run() = %v, want ErrNilSequence", err)
	}
}

func TestRunnerRunRejectsNegativeScheduleDelay(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(
		t,
		newTestClock(),
		WithSchedule(negativeDelaySchedule{}),
		WithInitialProbe(false),
	)

	err := runner.Run(context.Background())

	if !errors.Is(err, ErrInvalidScheduleDelay) {
		t.Fatalf("Run() = %v, want ErrInvalidScheduleDelay", err)
	}
}

func TestRunnerRunNilReceiver(t *testing.T) {
	t.Parallel()

	var runner *Runner
	err := runner.Run(context.Background())

	if !errors.Is(err, ErrNilRunner) {
		t.Fatalf("Run() = %v, want ErrNilRunner", err)
	}
}

func TestRunnerRunPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(t, newTestClock())

	defer func() {
		recovered := recover()
		if recovered != "healthprobe: nil context" {
			t.Fatalf("Run(nil) panic = %v, want healthprobe: nil context", recovered)
		}
	}()

	_ = runner.Run(nil)
}
