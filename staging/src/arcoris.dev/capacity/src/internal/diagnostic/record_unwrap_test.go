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

func TestRecordUnwrapPreservesSentinelAndCause(t *testing.T) {
	sentinel := errors.New("sentinel")
	cause := errors.New("cause")
	record := WrapRecord("", sentinel, "", "", cause)

	err := record.Unwrap()

	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is(err, sentinel) = false")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(err, cause) = false")
	}
}

func TestRecordUnwrapSingleErrors(t *testing.T) {
	sentinel := errors.New("sentinel")
	cause := errors.New("cause")

	if got := NewRecord("", sentinel, "", "").Unwrap(); got != sentinel {
		t.Fatalf("Record.Unwrap() = %v, want sentinel", got)
	}

	if got := CauseRecord[string](cause).Unwrap(); got != cause {
		t.Fatalf("Record.Unwrap() = %v, want cause", got)
	}

	if got := (Record[string]{}).Unwrap(); got != nil {
		t.Fatalf("Record.Unwrap() = %v, want nil", got)
	}
}
