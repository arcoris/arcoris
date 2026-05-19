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

package admission

import "testing"

func TestResultIsValidRejectsInvalidGrantShape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		result Result[string, NoMetadata]
	}{
		{
			name: "owned without grant",
			result: resultWith[string, NoMetadata](
				Grant(ReasonAdmitted),
				noneString(),
				noneMetadata(),
			),
		},
		{
			name: "denied with grant",
			result: resultWith(
				Deny(ReasonCapacityExhausted),
				someString("lease"),
				noneMetadata(),
			),
		},
		{
			name: "committed with grant",
			result: resultWith(
				Commit(ReasonAdmitted),
				someString("token"),
				noneMetadata(),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.result.IsValid() {
				t.Fatalf("%s result should be invalid", tt.name)
			}
		})
	}
}
