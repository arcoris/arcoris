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

package health

import "testing"

func TestTargetString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target Target
		want   string
	}{
		{TargetUnknown, "unknown"},
		{TargetStartup, "startup"},
		{TargetLive, "live"},
		{TargetReady, "ready"},
		{Target(99), "invalid"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.want, func(t *testing.T) {
			t.Parallel()

			if got := test.target.String(); got != test.want {
				t.Fatalf("String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTargetClassification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target   Target
		valid    bool
		concrete bool
		startup  bool
		live     bool
		ready    bool
	}{
		{TargetUnknown, true, false, false, false, false},
		{TargetStartup, true, true, true, false, false},
		{TargetLive, true, true, false, true, false},
		{TargetReady, true, true, false, false, true},
		{Target(99), false, false, false, false, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.target.String(), func(t *testing.T) {
			t.Parallel()

			if got := test.target.IsValid(); got != test.valid {
				t.Fatalf("IsValid() = %v, want %v", got, test.valid)
			}
			if got := test.target.IsConcrete(); got != test.concrete {
				t.Fatalf("IsConcrete() = %v, want %v", got, test.concrete)
			}
			if got := test.target.IsStartup(); got != test.startup {
				t.Fatalf("IsStartup() = %v, want %v", got, test.startup)
			}
			if got := test.target.IsLive(); got != test.live {
				t.Fatalf("IsLive() = %v, want %v", got, test.live)
			}
			if got := test.target.IsReady(); got != test.ready {
				t.Fatalf("IsReady() = %v, want %v", got, test.ready)
			}
		})
	}
}

func TestConcreteTargetsReturnsCallerOwnedOrder(t *testing.T) {
	t.Parallel()

	targets := ConcreteTargets()
	want := []Target{TargetStartup, TargetLive, TargetReady}
	if len(targets) != len(want) {
		t.Fatalf("len = %d, want %d", len(targets), len(want))
	}
	for i := range want {
		if targets[i] != want[i] {
			t.Fatalf("target[%d] = %s, want %s", i, targets[i], want[i])
		}
	}

	targets[0] = TargetReady
	if got := ConcreteTargets()[0]; got != TargetStartup {
		t.Fatalf("ConcreteTargets()[0] = %s, want startup", got)
	}
}
