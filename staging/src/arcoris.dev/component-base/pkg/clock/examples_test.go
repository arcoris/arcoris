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

package clock

import (
	"fmt"
	"time"
)

func exampleStartTime() time.Time {
	return time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
}

func observedAge(c PassiveClock, startedAt time.Time) time.Duration {
	return c.Since(startedAt)
}

// ExamplePassiveClock shows the narrow read-only clock contract.
//
// Components that only need Now or Since should accept PassiveClock instead of
// Clock. This keeps read-only code independent from timers, tickers, sleeps, and
// runtime loop ownership.
func ExamplePassiveClock() {
	clock := NewFakeClock(exampleStartTime())

	startedAt := clock.Now()
	clock.Step(250 * time.Millisecond)

	fmt.Println(observedAge(clock, startedAt))

	// Output:
	// 250ms
}

// ExampleClock_After shows a simple one-shot wait.
//
// After is useful when the caller does not need to stop or reset the wait. For
// cancelable or resettable waits, use NewTimer instead.
func ExampleClock_After() {
	clock := NewFakeClock(exampleStartTime())

	ch := clock.After(10 * time.Second)

	clock.Step(9 * time.Second)
	fmt.Println("before deadline")

	clock.Step(time.Second)
	fmt.Println((<-ch).Format(time.RFC3339))

	// Output:
	// before deadline
	// 2026-01-02T03:04:15Z
}

// ExampleClock_NewTimer shows a resettable one-shot timer.
//
// Timers are appropriate for dispatch timeouts, retry delays, controller
// cooldowns, queue wake-up deadlines, lease checks, and other one-shot waits
// that need explicit lifecycle ownership.
func ExampleClock_NewTimer() {
	clock := NewFakeClock(exampleStartTime())

	timer := clock.NewTimer(10 * time.Second)

	timer.Reset(3 * time.Second)
	clock.Step(3 * time.Second)

	fmt.Println((<-timer.C()).Format(time.RFC3339))

	// Output:
	// 2026-01-02T03:04:08Z
}

// ExampleClock_NewTicker shows a periodic fake ticker.
//
// Fake tickers are advanced explicitly by the owning FakeClock. This makes
// controller-loop tests deterministic and avoids real sleeps.
func ExampleClock_NewTicker() {
	clock := NewFakeClock(exampleStartTime())

	ticker := clock.NewTicker(5 * time.Second)
	defer ticker.Stop()

	clock.Step(5 * time.Second)
	fmt.Println((<-ticker.C()).Format(time.RFC3339))

	clock.Step(5 * time.Second)
	fmt.Println((<-ticker.C()).Format(time.RFC3339))

	// Output:
	// 2026-01-02T03:04:10Z
	// 2026-01-02T03:04:15Z
}

// ExampleTicker_Reset shows how a component can change loop cadence without
// rebuilding surrounding runtime state.
//
// Reset schedules the next tick relative to the clock's current time.
func ExampleTicker_Reset() {
	clock := NewFakeClock(exampleStartTime())

	ticker := clock.NewTicker(10 * time.Second)
	defer ticker.Stop()

	clock.Step(4 * time.Second)
	ticker.Reset(2 * time.Second)

	clock.Step(2 * time.Second)

	fmt.Println((<-ticker.C()).Format(time.RFC3339))

	// Output:
	// 2026-01-02T03:04:11Z
}

// ExampleFakeClock_Sleep shows deterministic Sleep behavior.
//
// FakeClock.Sleep blocks until another goroutine advances the same fake clock
// far enough to release the sleeper.
func ExampleFakeClock_Sleep() {
	clock := NewFakeClock(exampleStartTime())

	done := make(chan struct{})

	go func() {
		clock.Sleep(5 * time.Second)
		close(done)
	}()

	for !clock.HasWaiters() {
		// Spin until the goroutine has registered its fake-time waiter.
		//
		// Production code should not spin like this. This example uses the loop
		// only to keep fake-time behavior deterministic without real sleeps.
	}

	clock.Step(5 * time.Second)

	<-done
	fmt.Println("released")

	// Output:
	// released
}
