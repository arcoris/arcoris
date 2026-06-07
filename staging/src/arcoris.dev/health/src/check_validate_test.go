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

package health

import (
	"errors"
	"testing"
)

func TestValidateCheckerRejectsNilTypedNilAndInvalidNames(t *testing.T) {
	t.Parallel()

	var typed *typedNilChecker

	tests := []struct {
		name    string
		checker Checker
		want    error
	}{
		{name: "nil", checker: nil, want: ErrNilChecker},
		{name: "typed nil", checker: typed, want: ErrNilChecker},
		{name: "empty name", checker: checkerFunc{name: ""}, want: ErrEmptyCheckName},
		{name: "invalid name", checker: checkerFunc{name: "bad-name"}, want: ErrInvalidCheckName},
		{name: "valid", checker: mustCheck(t, "storage", Healthy("storage")), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateChecker(tc.checker)
			if !errors.Is(err, tc.want) {
				t.Fatalf("ValidateChecker() = %v, want %v", err, tc.want)
			}
		})
	}
}
