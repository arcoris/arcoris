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
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestStringOperation verifies string-only field operators and mismatch
// guardrails.
func TestStringOperation(t *testing.T) {
	actual := value.StringValue("api-server")

	if !stringOperation(actual, value.StringValue("api"), OperatorHasPrefix) {
		t.Fatal("HasPrefix failed")
	}
	if !stringOperation(actual, value.StringValue("server"), OperatorHasSuffix) {
		t.Fatal("HasSuffix failed")
	}
	if !stringOperation(actual, value.StringValue("-"), OperatorContains) {
		t.Fatal("Contains failed")
	}
	if stringOperation(value.Int64Value(1), value.StringValue("1"), OperatorContains) {
		t.Fatal("non-string actual matched")
	}
	if stringOperation(actual, value.Int64Value(1), OperatorContains) {
		t.Fatal("non-string literal matched")
	}
	if stringOperation(actual, value.StringValue("api"), OperatorEquals) {
		t.Fatal("unsupported string operator matched")
	}
}
