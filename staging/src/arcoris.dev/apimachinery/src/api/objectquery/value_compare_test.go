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

// TestValueEqualAndValueInUseCanonicalKeys verifies literal equality is shared
// across scalar and membership helpers.
func TestValueEqualAndValueInUseCanonicalKeys(t *testing.T) {
	if !valueEqual(value.StringValue("api"), value.StringValue("api")) {
		t.Fatal("equal strings were not equal")
	}
	if valueEqual(value.StringValue("1"), value.Int64Value(1)) {
		t.Fatal("different kinds compared equal")
	}
	if !valueIn(value.StringValue("api"), []value.Value{value.StringValue("web"), value.StringValue("api")}) {
		t.Fatal("valueIn did not find expected string")
	}
}
