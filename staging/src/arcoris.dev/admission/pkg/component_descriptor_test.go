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

func TestComponentDescriptorIsValid(t *testing.T) {
	t.Parallel()

	validID := MustComponentID("resilience.bulkhead")
	validCapabilities := NewCapabilitySet(
		CapabilityAdmit,
		CapabilityDeny,
		CapabilityEffectOwned,
		CapabilityEffectNone,
	)

	tests := []struct {
		name       string
		descriptor ComponentDescriptor
		want       bool
	}{
		{
			name: "valid descriptor",
			descriptor: ComponentDescriptor{
				ID:           validID,
				Kind:         KindBulkhead,
				Capabilities: validCapabilities,
			},
			want: true,
		},
		{
			name: "valid descriptor with unspecified capabilities",
			descriptor: ComponentDescriptor{
				ID:   validID,
				Kind: KindBulkhead,
			},
			want: true,
		},
		{
			name: "valid descriptor with syntactically valid unknown kind",
			descriptor: ComponentDescriptor{
				ID:           "custom.component",
				Kind:         "custom_kind",
				Capabilities: validCapabilities,
			},
			want: true,
		},
		{
			name: "invalid id",
			descriptor: ComponentDescriptor{
				ID:   "bad/id",
				Kind: KindBulkhead,
			},
		},
		{
			name: "invalid kind",
			descriptor: ComponentDescriptor{
				ID:   validID,
				Kind: "bad-kind",
			},
		},
		{
			name: "invalid capabilities",
			descriptor: ComponentDescriptor{
				ID:           validID,
				Kind:         KindBulkhead,
				Capabilities: CapabilitySet(1 << 15),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.descriptor.IsValid(); got != tt.want {
				t.Fatalf("IsValid = %v, want %v", got, tt.want)
			}
		})
	}
}
