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

package objectenqueue

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestIsNilInterface(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{name: "nil", in: nil, want: true},
		{name: "typed nil pointer", in: (*recordingQueue)(nil), want: true},
		{name: "typed nil func", in: MapperFunc(nil), want: true},
		{name: "non-nil pointer", in: &recordingQueue{}, want: false},
		{name: "non-nil func", in: MapperFunc(func() MapperFunc {
			return func(objectstore.Change, EmitFunc) error { return nil }
		}()), want: false},
		{name: "value", in: struct{}{}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNilInterface(tt.in)
			if got != tt.want {
				t.Fatalf("isNilInterface(%T) = %v; want %v", tt.in, got, tt.want)
			}
		})
	}
}
