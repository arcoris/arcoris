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

// RecordOriginal records one original, non-retry attempt.
//
// RecordOriginal contributes to the retry budget for the current fixed window.
// If the configured clock has reached the next window, the limiter rotates before
// recording the original attempt.
func (l *Limiter) RecordOriginal() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.cfg.clock.Now()
	l.rotateLocked(now)
	l.original = saturatingInc(l.original)
	l.publishLocked()
}
