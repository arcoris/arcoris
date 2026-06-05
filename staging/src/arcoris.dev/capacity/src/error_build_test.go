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

package capacity

import (
	"errors"
	"testing"
)

func TestErrorAtBuildsStructuredDiagnostic(t *testing.T) {
	t.Parallel()

	err := errorAt("path", ErrInvalidVector, ErrorReasonInvalidVector, "detail")
	if !errors.Is(err, ErrInvalidVector) {
		t.Fatalf("errorAt() = %v, want ErrInvalidVector", err)
	}

	capacityErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("errorAt() = %T, want *Error", err)
	}
	if capacityErr.Path != "path" || capacityErr.Reason != ErrorReasonInvalidVector || capacityErr.Detail != "detail" {
		t.Fatalf("record = %+v", capacityErr.Record)
	}
}
