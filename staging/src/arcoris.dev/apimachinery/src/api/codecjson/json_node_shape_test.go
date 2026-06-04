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
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestRequireObjectAcceptsObject(t *testing.T) {
	err := requireObject(rootPath(), jsonNode{kind: jsonKindObject}, "must be object")
	requireNoError(t, err)
}

func TestRequireObjectRejectsNonObject(t *testing.T) {
	err := requireObject(rootPath(), jsonNode{kind: jsonKindArray}, "must be object")

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, codec.ErrInvalidDocument)
}

func TestExpectStringAcceptsString(t *testing.T) {
	got, err := expectString(rootPath(), jsonNode{kind: jsonKindString, stringValue: "value"}, "must be string")
	requireNoError(t, err)

	if got != "value" {
		t.Fatalf("string = %q", got)
	}
}

func TestExpectStringRejectsNonString(t *testing.T) {
	_, err := expectString(rootPath(), jsonNode{kind: jsonKindNumber, numberText: "1"}, "must be string")

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireErrorIs(t, err, codec.ErrInvalidDocument)
}

func TestRejectUnknownMembers(t *testing.T) {
	node := jsonNode{kind: jsonKindObject, members: []jsonMember{{name: "unexpected"}}}

	err := rejectUnknownMembers(rootPath(), node, func(name string) bool {
		return name == "known"
	}, "test document")

	requireErrorIs(t, err, ErrInvalidEnvelope)
	requireCodecJSONError(t, err, "$.unexpected", ErrorReasonInvalidEnvelope)
	requireDetailContains(t, err, "test document")
}
