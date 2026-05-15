/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package fixedwindow

import (
	"errors"
	"testing"
	"time"
)

func TestValidateWindow(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want error
	}{
		{name: "valid", d: time.Second},
		{name: "zero", d: 0, want: ErrInvalidWindow},
		{name: "negative", d: -time.Second, want: ErrInvalidWindow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWindow(tt.d)
			if !errors.Is(err, tt.want) {
				t.Fatalf("validateWindow() error = %v, want %v", err, tt.want)
			}
		})
	}
}
