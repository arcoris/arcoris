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

package metagrammar

import "testing"

func TestASCIIHelpers(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		control    bool
		whitespace bool
		unsafe     bool
	}{
		{name: "plain", value: "worker-1"},
		{name: "control", value: "worker\n1", control: true, whitespace: true, unsafe: true},
		{name: "space", value: "worker 1", whitespace: true},
		{name: "slash", value: "worker/1", unsafe: true},
		{name: "backslash", value: `worker\1`, unsafe: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasControl(tt.value); got != tt.control {
				t.Fatalf("HasControl() = %v, want %v", got, tt.control)
			}
			if got := HasWhitespace(tt.value); got != tt.whitespace {
				t.Fatalf("HasWhitespace() = %v, want %v", got, tt.whitespace)
			}
			if got := HasUnsafeScalarChar(tt.value); got != tt.unsafe {
				t.Fatalf("HasUnsafeScalarChar() = %v, want %v", got, tt.unsafe)
			}
		})
	}
}
