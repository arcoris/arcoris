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

func TestWithObserverIgnoresNil(t *testing.T) {
	t.Parallel()

	cfg := newControllerConfig(WithObserver(nil))
	if len(cfg.observers) != 0 {
		t.Fatalf("observers len = %d, want 0", len(cfg.observers))
	}
}

func TestWithObserverAppendsOneObserver(t *testing.T) {
	t.Parallel()

	observer := ObserverFunc(func(Transition) {})
	cfg := newControllerConfig(WithObserver(observer))
	if len(cfg.observers) != 1 || cfg.observers[0] == nil {
		t.Fatalf("observers = %v, want one observer", cfg.observers)
	}
}

func TestWithObserversIgnoresNilEntries(t *testing.T) {
	t.Parallel()

	observer := ObserverFunc(func(Transition) {})
	cfg := newControllerConfig(WithObservers(nil, observer, nil))
	if len(cfg.observers) != 1 || cfg.observers[0] == nil {
		t.Fatalf("observers = %v, want only non-nil observer", cfg.observers)
	}
}

func TestWithObserversPreservesOrder(t *testing.T) {
	t.Parallel()

	// Observer ordering is a diagnostics and integration contract: observers see
	// the same committed transition in configuration order.
	var order []string
	first := ObserverFunc(func(Transition) { order = append(order, "first") })
	second := ObserverFunc(func(Transition) { order = append(order, "second") })
	third := ObserverFunc(func(Transition) { order = append(order, "third") })
	cfg := newControllerConfig(WithObserver(first), WithObservers(second, third))

	notifyObservers(cfg.observers, Transition{})
	assertDeepEqual(t, order, []string{"first", "second", "third"})
}
