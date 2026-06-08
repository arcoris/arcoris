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

package liveconfigtest

import "testing"

func TestEqualConfig(t *testing.T) {
	base := NewConfig()

	tests := []struct {
		name   string
		mutate func(*Config)
		want   bool
	}{
		{
			name: "equal",
			want: true,
		},
		{
			name:   "different name",
			mutate: func(cfg *Config) { cfg.Name = "other" },
		},
		{
			name:   "different limit",
			mutate: func(cfg *Config) { cfg.Limits[0] = 99 },
		},
		{
			name:   "different label",
			mutate: func(cfg *Config) { cfg.Labels["env"] = "prod" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CloneConfig(base)
			if tt.mutate != nil {
				tt.mutate(&got)
			}
			if EqualConfig(base, got) != tt.want {
				t.Fatalf("EqualConfig() = %v, want %v", EqualConfig(base, got), tt.want)
			}
		})
	}
}
