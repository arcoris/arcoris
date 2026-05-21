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

package retry

import (
	"context"
	"time"

	"arcoris.dev/chrono/clock"
)

// waitDelay waits for delay to elapse according to c or for ctx to stop.
//
// waitDelay is a retry-owned waiting primitive. It exists so retry execution can
// wait between attempts through an injected clock instead of using time.Sleep or
// time.Timer directly. This keeps retry loops deterministic with fake clocks in
// tests.
//
// waitDelay does not compute delay values, classify operation errors, update
// attempts, emit observer events, or create retry outcomes. It only owns the
// cancellable delay between two retry attempts.
//
// Semantics:
//
//   - delay < 0 is a programming error and panics;
//   - delay == 0 performs a retry-owned context pre-check and returns
//     immediately;
//   - delay > 0 creates a clock-owned timer;
//   - context stop before or during the wait is returned as ErrInterrupted;
//   - the timer is stopped when context wins the race.
//
// The returned context interruption is retry-owned. Raw context errors returned
// by operations must remain operation-owned unless retry itself observes the
// context stop at this boundary.
func waitDelay(ctx context.Context, c clock.Clock, d time.Duration) error {
	requireContext(ctx)
	requireClock(c)
	requireDelay(d, true)

	if err := contextStopError(ctx); err != nil {
		return err
	}
	if d == 0 {
		return nil
	}

	timer := c.NewTimer(d)

	select {
	case <-timer.C():
		return nil

	case <-ctx.Done():
		timer.Stop()
		return contextStopError(ctx)
	}
}
