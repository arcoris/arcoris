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

package admissioncatalog

import (
	"testing"

	"arcoris.dev/admission"
)

func TestReasonDescriptorIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		descriptor ReasonDescriptor
		want       bool
	}{
		{
			name: "built in reason with capabilities",
			descriptor: ReasonDescriptor{
				Reason: admission.ReasonAdmitted,
				Capabilities: NewCapabilitySet(
					CapabilityAdmit,
					CapabilityEffectNone,
				),
			},
			want: true,
		},
		{
			name: "custom open world reason",
			descriptor: ReasonDescriptor{
				Reason:       "custom_backoff",
				Capabilities: NewCapabilitySet(CapabilityDefer),
			},
			want: true,
		},
		{
			name: "zero capabilities are unspecified",
			descriptor: ReasonDescriptor{
				Reason: admission.ReasonDenied,
			},
			want: true,
		},
		{
			name: "invalid reason",
			descriptor: ReasonDescriptor{
				Reason:       "bad-reason",
				Capabilities: NewCapabilitySet(CapabilityDeny),
			},
			want: false,
		},
		{
			name: "invalid capabilities",
			descriptor: ReasonDescriptor{
				Reason:       admission.ReasonDenied,
				Capabilities: CapabilitySet(1 << 15),
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.descriptor.IsValid(); got != test.want {
				t.Fatalf("IsValid = %v, want %v", got, test.want)
			}
		})
	}
}
