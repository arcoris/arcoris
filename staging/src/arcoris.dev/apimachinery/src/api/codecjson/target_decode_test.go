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
)

func TestDecodeTargetFromUsesSharedDocumentParser(t *testing.T) {
	got, err := decodeTargetFrom(
		strings.NewReader(`"value"`),
		newTestCodec(t).decode,
		func(_ jsonPath, node jsonNode, _ resolvedDecodeConfig) (string, error) {
			return node.stringValue, nil
		},
	)
	requireNoError(t, err)

	if got != "value" {
		t.Fatalf("decoded = %q", got)
	}
}

func TestDecodeTargetFromReturnsZeroOnParserError(t *testing.T) {
	got, err := decodeTargetFrom(
		strings.NewReader(`"value" true`),
		newTestCodec(t).decode,
		func(_ jsonPath, _ jsonNode, _ resolvedDecodeConfig) (string, error) {
			return "unreachable", nil
		},
	)

	requireErrorIs(t, err, ErrTrailingData)
	if got != "" {
		t.Fatalf("decoded = %q; want zero value", got)
	}
}
