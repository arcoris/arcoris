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

func TestNewRecord(t *testing.T) {
	sentinel := errors.New("sentinel")

	record := NewRecord("value.path", sentinel, "invalid_value", "value is invalid")

	if record.Path != "value.path" {
		t.Fatalf("Path = %q, want value.path", record.Path)
	}
	if record.Err != sentinel {
		t.Fatalf("Err = %v, want sentinel", record.Err)
	}
	if record.Reason != "invalid_value" {
		t.Fatalf("Reason = %q, want invalid_value", record.Reason)
	}
	if record.Detail != "value is invalid" {
		t.Fatalf("Detail = %q, want value is invalid", record.Detail)
	}
	if record.Cause != nil {
		t.Fatalf("Cause = %v, want nil", record.Cause)
	}
}

func TestWrapRecord(t *testing.T) {
	sentinel := errors.New("sentinel")
	cause := errors.New("cause")

	record := WrapRecord("value.path", sentinel, "invalid_value", "value is invalid", cause)

	if record.Cause != cause {
		t.Fatalf("Cause = %v, want cause", record.Cause)
	}
}

func TestCauseRecord(t *testing.T) {
	cause := errors.New("cause")
	record := CauseRecord[string](cause)

	if record.Cause != cause {
		t.Fatalf("Cause = %v, want cause", record.Cause)
	}
	if record.Err != nil {
		t.Fatalf("Err = %v, want nil", record.Err)
	}
}

func TestJoinRecord(t *testing.T) {
	sentinel := errors.New("sentinel")
	cause := errors.New("cause")

	record := JoinRecord[string](sentinel, cause)

	if record.Err != sentinel {
		t.Fatalf("Err = %v, want sentinel", record.Err)
	}
	if record.Cause != cause {
		t.Fatalf("Cause = %v, want cause", record.Cause)
	}
}
