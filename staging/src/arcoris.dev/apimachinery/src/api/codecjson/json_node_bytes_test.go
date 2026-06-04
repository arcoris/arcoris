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

func TestJSONNodeBytes(t *testing.T) {
	data, err := jsonNodeBytes(jsonNode{kind: jsonKindString, stringValue: "<tag>"}, false)
	requireNoError(t, err)

	if string(data) != `"<tag>"` {
		t.Fatalf("bytes = %s", data)
	}
}

func TestJSONNodeBytesEscapesHTML(t *testing.T) {
	data, err := jsonNodeBytes(jsonNode{kind: jsonKindString, stringValue: "<tag>"}, true)
	requireNoError(t, err)

	if string(data) != `"\u003ctag\u003e"` {
		t.Fatalf("bytes = %s", data)
	}
}
