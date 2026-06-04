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

package codecjson

import "testing"

// TestJSONNodeMemberLookup covers private ordered object member helpers.
func TestJSONNodeMemberLookup(t *testing.T) {
	node := jsonNode{
		kind: jsonKindObject,
		members: []jsonMember{
			{name: "a", value: jsonNode{kind: jsonKindString, stringValue: "one"}},
			{name: "b", value: jsonNode{kind: jsonKindString, stringValue: "two"}},
		},
	}

	got, ok := node.member("b")
	if !ok {
		t.Fatalf("member b not found")
	}
	if got.stringValue != "two" {
		t.Fatalf("member b = %q; want two", got.stringValue)
	}
	if node.hasMember("missing") {
		t.Fatalf("missing member reported as present")
	}
}
