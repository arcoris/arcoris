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
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/clock"
)

func TestWithClock(t *testing.T) {
	fake := clock.NewFakeClock(time.Unix(10, 0))

	config := configOf(WithClock(fake))

	if config.clock != fake {
		t.Fatalf("configured clock was not stored")
	}
}

func TestWithClockPanicsOnNilClock(t *testing.T) {
	expectPanic(t, panicNilClock, func() {
		_ = WithClock(nil)
	})
}
