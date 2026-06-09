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

package objectlifecycle

import "testing"

func TestEffectStringAndValidity(t *testing.T) {
	tests := []struct {
		name  string
		eff   Effect
		text  string
		valid bool
	}{
		{name: "none", eff: EffectNone, text: "none", valid: true},
		{name: "found", eff: EffectFound, text: "found", valid: true},
		{name: "created", eff: EffectCreated, text: "created", valid: true},
		{name: "updated", eff: EffectUpdated, text: "updated", valid: true},
		{name: "unchanged", eff: EffectUnchanged, text: "unchanged", valid: true},
		{name: "deleted", eff: EffectDeleted, text: "deleted", valid: true},
		{name: "unknown", eff: Effect(99), text: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eff.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.eff.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
