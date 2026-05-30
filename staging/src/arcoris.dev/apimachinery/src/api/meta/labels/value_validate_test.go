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

package labels

import (
	"errors"
	"testing"
)

func TestValueValidate(t *testing.T) {
	valid := []Value{"", "worker", "worker-1", "worker_1", "worker.1", "1"}
	for _, value := range valid {
		t.Run("valid/"+value.String(), func(t *testing.T) {
			requireNoError(t, value.Validate())
		})
	}

	invalid := []Value{
		"worker value",
		"worker\nvalue",
		"-worker",
		"worker-",
		"_worker",
		"worker_",
		".worker",
		"worker.",
		".",
		"---",
		"___",
	}

	for _, value := range invalid {
		t.Run("invalid/"+value.String(), func(t *testing.T) {
			requireErrorIs(t, value.Validate(), ErrInvalidValue)
		})
	}
}

func TestValueValidateStructuredError(t *testing.T) {
	err := Value("-worker").Validate()
	requireErrorIs(t, err, ErrInvalidValue)

	var labelErr *Error
	if !errors.As(err, &labelErr) {
		t.Fatalf("errors.As(%T) = false", labelErr)
	}
	if labelErr.Path != "label.value" {
		t.Fatalf("Path = %q", labelErr.Path)
	}
	if labelErr.Reason != ErrorReasonInvalidEdge {
		t.Fatalf("Reason = %q", labelErr.Reason)
	}
	if labelErr.Detail == "" {
		t.Fatal("Detail is empty")
	}
}
