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
	sentinel := errors.New("invalid JSON")
	record := NewRecord("$.desired", sentinel, "invalid_value", "value is invalid")

	got := record.Format("codecjson")
	want := "codecjson: $.desired: invalid JSON: invalid_value: value is invalid"

	if got != want {
		t.Fatalf("Record.Format() = %q, want %q", got, want)
	}
}

func TestRecordUnwrapPreservesSentinelAndCause(t *testing.T) {
	sentinel := errors.New("sentinel")
	cause := errors.New("cause")
	record := WrapRecord("$", sentinel, "invalid_json", "invalid", cause)

	err := record.Unwrap()

	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is(err, sentinel) = false")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(err, cause) = false")
	}
}
