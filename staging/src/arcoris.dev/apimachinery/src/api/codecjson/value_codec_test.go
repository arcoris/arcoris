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
	"arcoris.dev/apimachinery/api/value"
	"testing"
)

// TestDecodeValueBytesRejectsDuplicateKey covers byte-slice value decoding.
func TestDecodeValueBytesRejectsDuplicateKey(t *testing.T) {
	_, err := newTestCodec(t).DecodeValue([]byte(`{"a":1,"a":2}`))

	requireErrorIs(t, err, ErrDuplicateKey)
}

// TestEncodeValueBytesUsesSharedWriter covers byte-slice value encoding.
func TestEncodeValueBytesUsesSharedWriter(t *testing.T) {
	data, err := newTestCodec(t).EncodeValue(value.StringValue("<tag>"))
	requireNoError(t, err)

	if got, want := string(data), `"<tag>"`; got != want {
		t.Fatalf("encoded value = %s; want %s", got, want)
	}
}
