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

func TestSetString(t *testing.T) {
	tests := []struct {
		name string
		set  Set
		want string
	}{
		{
			name: "empty",
			set:  EmptySet(),
			want: "{}",
		},
		{
			name: "single",
			set:  MustSet(setReplicasPath()),
			want: "{$.spec.replicas}",
		},
		{
			name: "multiple",
			set:  MustSet(setReplicasPath(), setImagePath(), setLabelPath()),
			want: `{$.metadata.labels["app"], $.spec.image, $.spec.replicas}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.set.String(), tt.want)
		})
	}
}
