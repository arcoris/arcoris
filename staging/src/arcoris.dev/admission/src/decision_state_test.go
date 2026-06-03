// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admission

import "testing"

func TestDecisionStateHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		decision      Decision
		admitted      bool
		denied        bool
		queued        bool
		deferred      bool
		terminal      bool
		sideEffect    bool
		requiresGrant bool
		allowsGrant   bool
	}{
		{
			name:        "plain admitted",
			decision:    Admit(ReasonAdmitted),
			admitted:    true,
			terminal:    true,
			allowsGrant: false,
		},
		{
			name:          "owned grant admitted",
			decision:      Grant(ReasonAdmitted),
			admitted:      true,
			terminal:      true,
			sideEffect:    true,
			requiresGrant: true,
			allowsGrant:   true,
		},
		{
			name:     "denied",
			decision: Deny(Reason("capacity_exhausted")),
			denied:   true,
			terminal: true,
		},
		{
			name:        "queued",
			decision:    Queue(ReasonQueued),
			queued:      true,
			sideEffect:  true,
			allowsGrant: true,
		},
		{
			name:     "deferred",
			decision: Defer(ReasonDeferred),
			deferred: true,
			terminal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.decision.IsAdmitted(); got != tt.admitted {
				t.Fatalf("IsAdmitted = %v, want %v", got, tt.admitted)
			}
			if got := tt.decision.IsDenied(); got != tt.denied {
				t.Fatalf("IsDenied = %v, want %v", got, tt.denied)
			}
			if got := tt.decision.IsQueued(); got != tt.queued {
				t.Fatalf("IsQueued = %v, want %v", got, tt.queued)
			}
			if got := tt.decision.IsDeferred(); got != tt.deferred {
				t.Fatalf("IsDeferred = %v, want %v", got, tt.deferred)
			}
			if got := tt.decision.IsTerminal(); got != tt.terminal {
				t.Fatalf("IsTerminal = %v, want %v", got, tt.terminal)
			}
			if got := tt.decision.HasSideEffect(); got != tt.sideEffect {
				t.Fatalf("HasSideEffect = %v, want %v", got, tt.sideEffect)
			}
			if got := tt.decision.RequiresGrant(); got != tt.requiresGrant {
				t.Fatalf("RequiresGrant = %v, want %v", got, tt.requiresGrant)
			}
			if got := tt.decision.AllowsGrant(); got != tt.allowsGrant {
				t.Fatalf("AllowsGrant = %v, want %v", got, tt.allowsGrant)
			}
		})
	}
}
