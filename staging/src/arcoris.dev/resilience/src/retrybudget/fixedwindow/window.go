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

import "time"

// rotateLocked starts a new accounting window when now is outside the current
// window.
//
// The caller must hold l.mu. Backward clock observations do not rotate the
// window. This keeps the limiter stable when a fake or wall clock moves
// backwards.
func (l *Limiter) rotateLocked(now time.Time) bool {
	if now.Before(l.windowStart) {
		return false
	}
	if now.Before(l.windowEndLocked()) {
		return false
	}

	l.windowStart = now
	l.original = 0
	l.retries = 0
	return true
}

// windowEndLocked returns the exclusive end of the current accounting window.
//
// The caller must hold l.mu or otherwise own a consistent view of l.windowStart.
func (l *Limiter) windowEndLocked() time.Time {
	return l.windowStart.Add(l.cfg.window)
}
