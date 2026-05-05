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

import "context"

// Run starts the fixed-interval probe loop and blocks until ctx is stopped.
//
// Run owns exactly one ticker loop for the Runner. Concurrent Run calls on the
// same Runner return ErrRunnerRunning. Context cancellation is treated as normal
// loop shutdown and returns nil.
//
// Run panics when ctx is nil. A nil context would create an unowned background
// loop and hide a caller wiring bug.
func (r *Runner) Run(ctx context.Context) error {
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

	ticker := r.clock.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C():
			r.runCycle(ctx)
		}
	}
}
