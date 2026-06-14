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

package objectsurface

import "testing"

func TestKindIsOwnable(t *testing.T) {
	kinds := Kinds()
	metadata := kinds.Metadata()

	tests := []struct {
		name string
		kind Kind
		want bool
	}{
		{name: "desired", kind: kinds.Desired(), want: true},
		{name: "observed", kind: kinds.Observed(), want: true},
		{name: "labels", kind: metadata.Labels(), want: true},
		{name: "annotations", kind: metadata.Annotations(), want: true},
		{name: "finalizers", kind: metadata.Finalizers(), want: false},
		{name: "owner references", kind: metadata.OwnerReferences(), want: false},
		{name: "empty", kind: "", want: false},
		{name: "unknown", kind: "metadata.name", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.IsOwnable(); got != tt.want {
				t.Fatalf("IsOwnable() = %v; want %v", got, tt.want)
			}
		})
	}
}
