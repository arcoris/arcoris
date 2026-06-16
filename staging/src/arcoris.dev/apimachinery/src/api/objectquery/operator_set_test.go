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

// TestOperatorSetIgnoresInvalidOperators verifies defensive construction does
// not accidentally enable unknown operators.
func TestOperatorSetIgnoresInvalidOperators(t *testing.T) {
	set := Operators(OperatorEquals, Operator(99))

	if !set.Supports(OperatorEquals) {
		t.Fatal("set does not support equals")
	}
	if set.Supports(Operator(99)) {
		t.Fatal("set supports invalid operator")
	}
}

// TestMetadataOperatorsRemainMetadataOnly verifies metadata maps do not gain
// field-only operators through the generalized operator enum.
func TestMetadataOperatorsRemainMetadataOnly(t *testing.T) {
	if !metadataOperators.Supports(OperatorIn) {
		t.Fatal("metadata operators do not support in")
	}
	if metadataOperators.Supports(OperatorContains) {
		t.Fatal("metadata operators support contains; want field-only")
	}
}
