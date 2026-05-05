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

package health

import (
	"errors"
	"strings"
	"testing"
)

func TestValidCheckName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want bool
	}{
		{name: "database", want: true},
		{name: "database_pool_1", want: true},
		{name: "", want: false},
		{name: "1database", want: false},
		{name: "_database", want: false},
		{name: "database_", want: false},
		{name: "database__pool", want: false},
		{name: "database-pool", want: false},
		{name: "Database", want: false},
		{name: strings.Repeat("a", maxCheckNameLength+1), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := ValidCheckName(tc.name); got != tc.want {
				t.Fatalf("ValidCheckName(%q) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestValidateCheckName(t *testing.T) {
	t.Parallel()

	if err := ValidateCheckName("storage"); err != nil {
		t.Fatalf("ValidateCheckName(valid) = %v, want nil", err)
	}

	if err := ValidateCheckName(""); !errors.Is(err, ErrEmptyCheckName) {
		t.Fatalf("ValidateCheckName(empty) = %v, want ErrEmptyCheckName", err)
	}

	if err := ValidateCheckName("bad-name"); !errors.Is(err, ErrInvalidCheckName) {
		t.Fatalf("ValidateCheckName(invalid) = %v, want ErrInvalidCheckName", err)
	}
}
