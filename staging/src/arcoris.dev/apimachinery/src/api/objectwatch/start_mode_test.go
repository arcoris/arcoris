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

package objectwatch

import "testing"

func TestStartModeStringAndValidity(t *testing.T) {
	tests := []struct {
		mode  StartMode
		text  string
		valid bool
	}{
		{mode: 0, text: "unknown"},
		{mode: StartAfterRevision, text: "afterRevision", valid: true},
		{mode: StartAtCurrent, text: "atCurrent", valid: true},
		{mode: StartMode(99), text: "unknown"},
	}

	for _, tt := range tests {
		if tt.mode.String() != tt.text {
			t.Fatalf("String(%d) = %q; want %q", tt.mode, tt.mode.String(), tt.text)
		}
		if tt.mode.IsValid() != tt.valid {
			t.Fatalf("IsValid(%d) = %v; want %v", tt.mode, tt.mode.IsValid(), tt.valid)
		}
	}
}
