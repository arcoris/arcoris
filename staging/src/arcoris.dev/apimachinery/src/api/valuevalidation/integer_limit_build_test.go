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

package valuevalidation

import "testing"

func TestIntegerLimitConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  integerLimits[int64]
		want integerLimits[int64]
	}{
		{
			name: "signed width",
			got:  signedWidthLimits(-8, 7),
			want: integerLimits[int64]{
				lower: integerBound[int64]{value: -8, set: true},
				upper: integerBound[int64]{value: 7, set: true},
			},
		},
		{
			name: "small signed descriptor",
			got: signedDescriptorLimits[int8](
				func() (int8, bool) { return -2, true },
				func() (int8, bool) { return 3, true },
			),
			want: integerLimits[int64]{
				lower: integerBound[int64]{value: -2, set: true},
				upper: integerBound[int64]{value: 3, set: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("limits = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestUnsignedLimitConstructors(t *testing.T) {
	got := unsignedDescriptorLimits[uint8](
		func() (uint8, bool) { return 1, true },
		func() (uint8, bool) { return 9, true },
	)
	want := integerLimits[uint64]{
		lower: integerBound[uint64]{value: 1, set: true},
		upper: integerBound[uint64]{value: 9, set: true},
	}

	if got != want {
		t.Fatalf("limits = %#v, want %#v", got, want)
	}
}
