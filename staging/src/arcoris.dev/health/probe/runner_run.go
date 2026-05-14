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

import (
	"context"
	"time"

	"arcoris.dev/chrono/delay"
)

// Run starts the schedule-driven probe loop and blocks until ctx is stopped.
//
// Run owns exactly one schedule sequence for the Runner. Concurrent Run calls
// on the same Runner return ErrRunnerRunning. Context cancellation is treated as
// normal loop shutdown and returns nil. Schedule exhaustion returns
// ErrExhaustedSchedule. Calling Run on a nil Runner returns ErrNilRunner.
//
// Run panics when ctx is nil. A nil context would create an unowned background
// loop and hide a caller wiring bug.
func (r *Runner) Run(ctx context.Context) error {
	if r == nil {
		return ErrNilRunner
	}
	if ctx == nil {
		panic("healthprobe: nil context")
	}
	if !r.running.CompareAndSwap(false, true) {
		return ErrRunnerRunning
	}
	defer r.running.Store(false)

	if r.initialProbe {
		r.runCycle(ctx)
		if ctx.Err() != nil {
			return nil
		}
	}

	if ctx.Err() != nil {
		return nil
	}

	sequence, err := newSequence(r.schedule)
	if err != nil {
		return err
	}

	for {
		wait, ok, err := nextDelay(ctx, sequence)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if !r.waitDelay(ctx, wait) {
			return nil
		}

		r.runCycle(ctx)
		if ctx.Err() != nil {
			return nil
		}
	}
}

func nextDelay(ctx context.Context, sequence delay.Sequence) (time.Duration, bool, error) {
	d, ok := sequence.Next()
	if !ok {
		if ctx.Err() != nil {
			return 0, false, nil
		}
		return 0, false, ErrExhaustedSchedule
	}
	if d < 0 {
		return 0, false, InvalidScheduleDelayError{Delay: d}
	}

	return d, true, nil
}

func (r *Runner) waitDelay(ctx context.Context, d time.Duration) bool {
	if d == 0 {
		return true
	}

	timer := r.clock.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C():
		return true
	}
}
