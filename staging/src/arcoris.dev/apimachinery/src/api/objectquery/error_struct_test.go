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

package objectquery

import (
	"strings"
	"testing"
)

func TestErrorNilBehavior(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q; want <nil>", got)
	}
	if got := err.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v; want nil", got)
	}
}

func TestErrorFormatsPackageName(t *testing.T) {
	err := errorf("query.labels", ErrInvalidQuery, ErrorReasonInvalidQuery, "bad query")
	if got := err.Error(); !strings.Contains(got, "objectquery") {
		t.Fatalf("Error() = %q; want package name", got)
	}
}
