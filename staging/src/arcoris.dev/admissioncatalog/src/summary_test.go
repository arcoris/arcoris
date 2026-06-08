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

import "testing"

func TestValidSummary(t *testing.T) {
	tests := []struct {
		name    string
		summary string
		want    bool
	}{
		{name: "empty", summary: "", want: true},
		{name: "plain", summary: "Static catalog documentation.", want: true},
		{name: "punctuation", summary: "Static, human-facing metadata: ok!", want: true},
		{name: "newline", summary: "dynamic\nlog", want: false},
		{name: "tab", summary: "dynamic\tlog", want: false},
		{name: "nul", summary: string(rune(0)), want: false},
		{name: "delete", summary: string(rune(0x7f)), want: false},
		{name: "invalid utf8", summary: string([]byte{0xff}), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validSummary(tt.summary); got != tt.want {
				t.Fatalf("validSummary(%q) = %v, want %v", tt.summary, got, tt.want)
			}
		})
	}
}
