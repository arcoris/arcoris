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

func TestCloneSetPathDetachesElements(t *testing.T) {
	original := setReplicasPath()
	cloned := cloneSetPath(original)

	cloned.elements[0] = testFieldElement("status")

	requireEqual(t, original.String(), "$.spec.replicas")
	requireEqual(t, cloned.String(), "$.status.replicas")
}

func TestCompareSetPathSlices(t *testing.T) {
	tests := []struct {
		name  string
		left  []Path
		right []Path
		want  int
	}{
		{
			name:  "equal",
			left:  []Path{setImagePath()},
			right: []Path{setImagePath()},
			want:  0,
		},
		{
			name:  "path order",
			left:  []Path{setImagePath()},
			right: []Path{setReplicasPath()},
			want:  -1,
		},
		{
			name:  "shorter",
			left:  []Path{setImagePath()},
			right: []Path{setImagePath(), setReplicasPath()},
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, compareSetPathSlices(tt.left, tt.right), tt.want)
		})
	}
}
