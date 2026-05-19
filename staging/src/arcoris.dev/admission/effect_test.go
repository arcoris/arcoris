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

func TestEffectIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		effect Effect
		want   bool
	}{
		{name: "unknown", effect: EffectUnknown, want: false},
		{name: "none", effect: EffectNone, want: true},
		{name: "committed", effect: EffectCommitted, want: true},
		{name: "owned", effect: EffectOwned, want: true},
		{name: "queued", effect: EffectQueued, want: true},
		{name: "undefined", effect: Effect(99), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.effect.IsValid(); got != tt.want {
				t.Fatalf("%v IsValid = %v, want %v", tt.effect, got, tt.want)
			}
		})
	}
}

func TestEffectGrantSemantics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		effect        Effect
		sideEffect    bool
		requiresGrant bool
		allowsGrant   bool
	}{
		{name: "none", effect: EffectNone},
		{name: "committed", effect: EffectCommitted, sideEffect: true},
		{
			name:          "owned",
			effect:        EffectOwned,
			sideEffect:    true,
			requiresGrant: true,
			allowsGrant:   true,
		},
		{
			name:        "queued",
			effect:      EffectQueued,
			sideEffect:  true,
			allowsGrant: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.effect.HasSideEffect(); got != tt.sideEffect {
				t.Fatalf("HasSideEffect = %v, want %v", got, tt.sideEffect)
			}
			if got := tt.effect.RequiresGrant(); got != tt.requiresGrant {
				t.Fatalf("RequiresGrant = %v, want %v", got, tt.requiresGrant)
			}
			if got := tt.effect.AllowsGrant(); got != tt.allowsGrant {
				t.Fatalf("AllowsGrant = %v, want %v", got, tt.allowsGrant)
			}
		})
	}
}

func TestEffectString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		effect Effect
		want   string
	}{
		{name: "unknown", effect: EffectUnknown, want: "unknown"},
		{name: "none", effect: EffectNone, want: "none"},
		{name: "committed", effect: EffectCommitted, want: "committed"},
		{name: "owned", effect: EffectOwned, want: "owned"},
		{name: "queued", effect: EffectQueued, want: "queued"},
		{name: "undefined", effect: Effect(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.effect.String(); got != tt.want {
				t.Fatalf("String = %q, want %q", got, tt.want)
			}
		})
	}
}
