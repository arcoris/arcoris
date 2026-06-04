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
	"arcoris.dev/apimachinery/api/value"
)

// TestDecodeObjectBytesRejectsTrailingData covers byte-slice object decoding.
func TestDecodeObjectBytesRejectsTrailingData(t *testing.T) {
	_, err := newTestCodec(t).DecodeObject([]byte(`{"desired":null} {}`))

	requireErrorIs(t, err, ErrTrailingData)
}

// TestEncodeObjectBytesUsesSharedWriter covers byte-slice object encoding.
func TestEncodeObjectBytesUsesSharedWriter(t *testing.T) {
	data, err := newTestCodec(t).EncodeObject(codec.Object{Desired: value.NullValue()})
	requireNoError(t, err)

	if got, want := string(data), `{"desired":null}`; got != want {
		t.Fatalf("encoded object = %s; want %s", got, want)
	}
}
