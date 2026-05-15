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

package fixedwindow

import "arcoris.dev/resilience/retrybudget"

// TryAdmitRetry atomically decides whether one retry attempt may be admitted.
//
// When admission succeeds, the retry attempt is recorded before the decision is
// returned. When admission is denied without a window rotation, the limiter state
// is unchanged and the returned decision carries the latest already-published
// snapshot.
func (l *Limiter) TryAdmitRetry() retrybudget.Decision {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.cfg.clock.Now()
	rotated := l.rotateLocked(now)
	allowed := allowedRetries(l.original, l.cfg.ratio, l.cfg.minRetries)

	if l.retries >= allowed {
		snap := l.published.Snapshot()
		if rotated {
			snap = l.publishLocked()
		}
		return retrybudget.Decision{
			Allowed:  false,
			Reason:   retrybudget.ReasonExhausted,
			Snapshot: snap,
		}
	}

	l.retries = saturatingInc(l.retries)
	return retrybudget.Decision{
		Allowed:  true,
		Reason:   retrybudget.ReasonAllowed,
		Snapshot: l.publishLocked(),
	}
}
