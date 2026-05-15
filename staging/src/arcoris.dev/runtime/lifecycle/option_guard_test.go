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

func TestWithGuardIgnoresNil(t *testing.T) {
	t.Parallel()

	cfg := newControllerConfig(WithGuard(nil))
	if len(cfg.guards) != 0 {
		t.Fatalf("guards len = %d, want 0", len(cfg.guards))
	}
}

func TestWithGuardAppendsOneGuard(t *testing.T) {
	t.Parallel()

	guard := TransitionGuardFunc(func(Transition) error { return nil })
	cfg := newControllerConfig(WithGuard(guard))
	if len(cfg.guards) != 1 || cfg.guards[0] == nil {
		t.Fatalf("guards = %v, want one guard", cfg.guards)
	}
}

func TestWithGuardsIgnoresNilEntries(t *testing.T) {
	t.Parallel()

	guard := TransitionGuardFunc(func(Transition) error { return nil })
	cfg := newControllerConfig(WithGuards(nil, guard, nil))
	if len(cfg.guards) != 1 || cfg.guards[0] == nil {
		t.Fatalf("guards = %v, want only non-nil guard", cfg.guards)
	}
}

func TestWithGuardsPreservesOrder(t *testing.T) {
	t.Parallel()

	// Guard ordering is controller semantics because the first rejection becomes
	// the returned domain cause.
	var order []string
	first := TransitionGuardFunc(func(Transition) error { order = append(order, "first"); return nil })
	second := TransitionGuardFunc(func(Transition) error { order = append(order, "second"); return nil })
	third := TransitionGuardFunc(func(Transition) error { order = append(order, "third"); return nil })
	cfg := newControllerConfig(WithGuard(first), WithGuards(second, third))

	if err := allowTransition(cfg.guards, Transition{}); err != nil {
		t.Fatalf("allowTransition = %v, want nil", err)
	}
	assertDeepEqual(t, order, []string{"first", "second", "third"})
}
