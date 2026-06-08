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

import (
	"errors"
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want error
	}{
		{
			name: "valid",
			cfg:  NewConfig(),
		},
		{
			name: "blank name",
			cfg:  InvalidNameConfig(),
			want: ErrInvalidName,
		},
		{
			name: "negative version",
			cfg:  InvalidVersionConfig(),
			want: ErrInvalidVersion,
		},
		{
			name: "zero timeout",
			cfg:  InvalidTimeoutConfig(),
			want: ErrInvalidTimeout,
		},
		{
			name: "negative limit",
			cfg:  InvalidLimitConfig(),
			want: ErrInvalidLimit,
		},
		{
			name: "positive timeout",
			cfg: func() Config {
				cfg := NewConfig()
				cfg.Timeout = time.Nanosecond
				return cfg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if !errors.Is(err, tt.want) {
				t.Fatalf("ValidateConfig() error = %v, want %v", err, tt.want)
			}
		})
	}
}
