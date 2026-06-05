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

package diagnostic

import (
	"errors"
	"testing"
)

func TestRecordFormat(t *testing.T) {
	sentinel := errors.New("invalid value")
	record := NewRecord("object.members[0]", sentinel, "invalid_member", "member is invalid")

	got := record.Format("value")
	want := "value: object.members[0]: invalid value: invalid_member: member is invalid"

	if got != want {
		t.Fatalf("Record.Format() = %q, want %q", got, want)
	}
}

func TestRecordFormatSkipsEmptyParts(t *testing.T) {
	record := NewRecord("", nil, "invalid_syntax", "")

	got := record.Format("fieldpath")
	want := "fieldpath: invalid_syntax"

	if got != want {
		t.Fatalf("Record.Format() = %q, want %q", got, want)
	}
}

func TestRecordFormatAcceptsPackageReasonTypes(t *testing.T) {
	type localReason string

	record := NewRecord("", nil, localReason("local_reason"), "")

	if got, want := record.Format("domain"), "domain: local_reason"; got != want {
		t.Fatalf("Record.Format() = %q, want %q", got, want)
	}
}
