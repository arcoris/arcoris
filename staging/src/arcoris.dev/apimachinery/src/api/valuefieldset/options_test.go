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

package valuefieldset

import "testing"

func TestOptionsNormalizedMaxDepth(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want int
	}{
		{name: "default", opts: Options{}, want: defaultMaxDepth},
		{name: "negative", opts: Options{MaxDepth: -1}, want: defaultMaxDepth},
		{name: "custom", opts: Options{MaxDepth: 7}, want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.normalizedMaxDepth(); got != tt.want {
				t.Fatalf("normalizedMaxDepth() = %d, want %d", got, tt.want)
			}
		})
	}
}
