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

package lifecycle

import "testing"

func TestOptions(t *testing.T) {
	t.Parallel()

	guard := TransitionGuardFunc(func(Transition) error { return nil })
	observer := ObserverFunc(func(Transition) {})

	config := newControllerConfig(
		nil,
		WithClock(nil),
		WithClock(testClock{now: testTime}),
		WithGuard(nil),
		WithGuard(guard),
		WithGuards(nil, guard),
		WithObserver(nil),
		WithObserver(observer),
		WithObservers(nil, observer),
	)

	if got := config.now(); !got.Equal(testTime) {
		t.Fatalf("config.now() = %v, want %v", got, testTime)
	}
	if len(config.guards) != 2 {
		t.Fatalf("guards len = %d, want 2", len(config.guards))
	}
	if len(config.observers) != 2 {
		t.Fatalf("observers len = %d, want 2", len(config.observers))
	}
}
