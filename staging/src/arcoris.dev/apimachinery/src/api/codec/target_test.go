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

package codec

import "testing"

func TestTargetString(t *testing.T) {
	if TargetObjectOwnership.String() != "object_ownership" {
		t.Fatalf("Target.String() = %q", TargetObjectOwnership.String())
	}
}

func TestTargetIsZero(t *testing.T) {
	if !Target("").IsZero() {
		t.Fatalf("zero target IsZero() = false")
	}
	if TargetValue.IsZero() {
		t.Fatalf("value target IsZero() = true")
	}
}
