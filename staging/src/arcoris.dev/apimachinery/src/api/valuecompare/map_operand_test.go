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

package valuecompare

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestMapOperandReportsPresence(t *testing.T) {
	members := map[string]value.Value{"env": value.StringValue("prod")}

	got := mapOperand(members, "env")
	text, _ := got.value.String()
	if !got.present || text != "prod" {
		t.Fatalf("mapOperand(existing) = %#v", got)
	}

	got = mapOperand(members, "missing")
	if got.present {
		t.Fatalf("mapOperand(missing).present = true")
	}
}
