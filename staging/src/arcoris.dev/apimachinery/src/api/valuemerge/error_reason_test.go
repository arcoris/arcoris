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

import "testing"

func TestErrorReasonValues(t *testing.T) {
	if ErrorReasonInvalidPath != "invalid_path" {
		t.Fatalf("invalid path reason = %q", ErrorReasonInvalidPath)
	}
	if ErrorReasonUnsupportedMerge != "unsupported_merge" {
		t.Fatalf("unsupported merge reason = %q", ErrorReasonUnsupportedMerge)
	}
	if ErrorReasonInvalidFieldName != "invalid_field_name" {
		t.Fatalf("invalid field name reason = %q", ErrorReasonInvalidFieldName)
	}
	if ErrorReasonInvalidMapKey != "invalid_map_key" {
		t.Fatalf("invalid map key reason = %q", ErrorReasonInvalidMapKey)
	}
}
