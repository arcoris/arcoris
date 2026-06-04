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

import (
	"bytes"
	"testing"
)

func TestEncodeJSONDocumentPreservesOrder(t *testing.T) {
	node := jsonNode{kind: jsonKindObject, members: []jsonMember{
		{name: "b", value: jsonNode{kind: jsonKindNumber, numberText: "1"}},
		{name: "a", value: jsonNode{kind: jsonKindNumber, numberText: "2"}},
	}}

	var buffer bytes.Buffer
	requireNoError(t, encodeJSONDocument(&buffer, node, nodeEncodeConfig{}))

	if got := buffer.String(); got != `{"b":1,"a":2}` {
		t.Fatalf("encoded = %s", got)
	}
}

func TestEncodeJSONDocumentPretty(t *testing.T) {
	node := jsonNode{kind: jsonKindObject, members: []jsonMember{
		{name: "a", value: jsonNode{kind: jsonKindNull}},
	}}

	var buffer bytes.Buffer
	requireNoError(t, encodeJSONDocument(&buffer, node, nodeEncodeConfig{pretty: true, indent: "  "}))

	if got := buffer.String(); got != "{\n  \"a\": null\n}" {
		t.Fatalf("encoded = %q", got)
	}
}
