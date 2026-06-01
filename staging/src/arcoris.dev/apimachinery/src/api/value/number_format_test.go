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

package value

import "testing"

func TestAppendPaddedUnsignedDecimal(t *testing.T) {
	tests := []struct {
		name  string
		value uint64
		width int
		want  string
	}{
		{name: "pads short value", value: 7, width: 3, want: "007"},
		{name: "keeps exact width", value: 12, width: 2, want: "12"},
		{name: "keeps wide value", value: 123, width: 2, want: "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(appendPaddedUnsignedDecimal(nil, tt.value, tt.width))

			requireEqual(t, got, tt.want)
		})
	}
}

func TestAppendPaddedSignedDecimal(t *testing.T) {
	tests := []struct {
		name  string
		value int
		width int
		want  string
	}{
		{name: "positive", value: 7, width: 4, want: "0007"},
		{name: "negative", value: -7, width: 4, want: "-007"},
		{name: "wide negative", value: -1234, width: 4, want: "-1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(appendPaddedSignedDecimal(nil, tt.value, tt.width))

			requireEqual(t, got, tt.want)
		})
	}
}
