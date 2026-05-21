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

func TestComponentKindDescriptorIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		descriptor ComponentKindDescriptor
		want       bool
	}{
		{
			name: "valid descriptor",
			descriptor: ComponentKindDescriptor{
				Kind:         "custom_guard",
				Capabilities: NewCapabilitySet(CapabilityAdmit, CapabilityDeny),
			},
			want: true,
		},
		{
			name: "valid zero capabilities",
			descriptor: ComponentKindDescriptor{
				Kind: "custom_guard",
			},
			want: true,
		},
		{
			name: "invalid kind",
			descriptor: ComponentKindDescriptor{
				Kind:         "bad-kind",
				Capabilities: NewCapabilitySet(CapabilityAdmit),
			},
		},
		{
			name: "invalid capabilities",
			descriptor: ComponentKindDescriptor{
				Kind:         "custom_guard",
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
