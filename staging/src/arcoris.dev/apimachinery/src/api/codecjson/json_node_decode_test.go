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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestDecodeJSONDocumentPreservesObjectOrder(t *testing.T) {
	node, err := decodeJSONDocument(strings.NewReader(`{"b":1,"a":2}`), newTestCodec(t).decode)
	requireNoError(t, err)

	if got := []string{node.members[0].name, node.members[1].name}; got[0] != "b" || got[1] != "a" {
		t.Fatalf("member order = %#v", got)
	}
}

func TestDecodeJSONDocumentRejectsDuplicateKey(t *testing.T) {
	_, err := decodeJSONDocument(strings.NewReader(`{"a":1,"a":2}`), newTestCodec(t).decode)

	requireErrorIs(t, err, ErrDuplicateKey)
	requireErrorIs(t, err, codec.ErrInvalidDocument)
	requireCodecJSONError(t, err, "$", ErrorReasonDuplicateKey)
}

func TestDecodeJSONDocumentRejectsTrailingData(t *testing.T) {
	_, err := decodeJSONDocument(strings.NewReader(`null null`), newTestCodec(t).decode)

	requireErrorIs(t, err, ErrTrailingData)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
}

func TestDecodeJSONDocumentRejectsMalformedJSON(t *testing.T) {
	_, err := decodeJSONDocument(strings.NewReader(`{"a":`), newTestCodec(t).decode)

	requireErrorIs(t, err, ErrInvalidJSON)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
}

func TestDecodeJSONDocumentRejectsMaxDepthExceeded(t *testing.T) {
	config := newTestCodec(t).decode
	config.maxDepth = 2

	_, err := decodeJSONDocument(strings.NewReader(`[[null]]`), config)

	requireErrorIs(t, err, codec.ErrDepthExceeded)
	requireCodecJSONError(t, err, "$[0][0]", ErrorReasonMaxDepthExceeded)
}
