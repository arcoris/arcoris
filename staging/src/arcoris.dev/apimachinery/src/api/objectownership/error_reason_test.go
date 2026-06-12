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

package objectownership

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	if ErrorReasonInvalidDocument != "invalid_document" {
		t.Fatalf("invalid document reason = %q", ErrorReasonInvalidDocument)
	}
	if ErrorReasonMissingVersion != "missing_version" {
		t.Fatalf("missing version reason = %q", ErrorReasonMissingVersion)
	}
	if ErrorReasonUnsupportedVersion != "unsupported_version" {
		t.Fatalf("unsupported version reason = %q", ErrorReasonUnsupportedVersion)
	}
	if ErrorReasonInvalidPath != "invalid_path" {
		t.Fatalf("invalid path reason = %q", ErrorReasonInvalidPath)
	}
	if ErrorReasonNotNormalized != "not_normalized" {
		t.Fatalf("not normalized reason = %q", ErrorReasonNotNormalized)
	}
}
