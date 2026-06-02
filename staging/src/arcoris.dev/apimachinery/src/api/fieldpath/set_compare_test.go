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

package fieldpath

import "testing"

func TestSetEqual(t *testing.T) {
	left := MustSet(setReplicasPath(), setImagePath())
	right := MustSet(setImagePath(), setReplicasPath(), setReplicasPath())

	requireEqual(t, left.Equal(right), true)
}

func TestSetCompare(t *testing.T) {
	tests := []struct {
		name  string
		left  Set
		right Set
		want  int
	}{
		{
			name:  "equal",
			left:  MustSet(setImagePath()),
			right: MustSet(setImagePath()),
			want:  0,
		},
		{
			name:  "path ordering",
			left:  MustSet(setImagePath()),
			right: MustSet(setReplicasPath()),
			want:  -1,
		},
		{
			name:  "shorter prefix",
			left:  MustSet(setImagePath()),
			right: MustSet(setImagePath(), setReplicasPath()),
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.left.Compare(tt.right), tt.want)
		})
	}
}
