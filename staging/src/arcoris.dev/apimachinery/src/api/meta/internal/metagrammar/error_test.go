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

package metagrammar

import (
	"strings"
	"testing"
)

func TestViolationError(t *testing.T) {
	err := violation(ReasonInvalidForm, "bad form")
	if got := err.Error(); !strings.Contains(got, "invalid_form") || !strings.Contains(got, "bad form") {
		t.Fatalf("Error() = %q", got)
	}
}

func TestReasonValues(t *testing.T) {
	want := map[Reason]string{
		ReasonEmptyValue:       "empty_value",
		ReasonInvalidLength:    "invalid_length",
		ReasonInvalidCharacter: "invalid_character",
		ReasonInvalidEdge:      "invalid_edge",
		ReasonInvalidForm:      "invalid_form",
	}
	for reason, text := range want {
		if string(reason) != text {
			t.Fatalf("%#v = %q, want %q", reason, reason, text)
		}
	}
}
