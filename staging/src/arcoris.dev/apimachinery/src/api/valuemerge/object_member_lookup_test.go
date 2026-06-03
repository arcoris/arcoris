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

package valuemerge

import (
	"testing"

	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/value"
)

func TestMemberLookupOperandPresent(t *testing.T) {
	lookup := newMemberLookup([]value.Member{
		member("name", str("api")),
	})

	got := lookup.Operand("name")
	if got.Absent() {
		t.Fatalf("operand is absent")
	}

	text, _ := got.Value().String()
	if text != "api" {
		t.Fatalf("value = %q; want api", text)
	}
}

func TestMemberLookupOperandAbsent(t *testing.T) {
	got := newMemberLookup(nil).Operand("missing")

	if got.Present() {
		t.Fatalf("operand is present")
	}
}

func TestAppendMemberSkipsAbsent(t *testing.T) {
	got := appendMember(nil, "name", valuepresence.Absent())

	if len(got) != 0 {
		t.Fatalf("members length = %d; want 0", len(got))
	}
}
