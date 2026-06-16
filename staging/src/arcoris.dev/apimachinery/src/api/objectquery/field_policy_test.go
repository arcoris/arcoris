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

import "testing"

// TestFieldPolicyConstants locks down zero-value policy choices because they
// are embedded in selectable field declarations.
func TestFieldPolicyConstants(t *testing.T) {
	if IndexNone != 0 {
		t.Fatalf("IndexNone = %d; want zero", IndexNone)
	}
	if MissingAbsent != 0 {
		t.Fatalf("MissingAbsent = %d; want zero", MissingAbsent)
	}
	if IndexEquality == IndexRange {
		t.Fatal("IndexEquality and IndexRange must remain distinct")
	}
}
