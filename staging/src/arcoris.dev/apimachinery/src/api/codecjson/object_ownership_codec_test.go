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
	"arcoris.dev/apimachinery/api/objectownership"
	"testing"
)

// TestDecodeObjectOwnershipBytesRejectsDuplicateKey covers byte-slice ownership decoding.
func TestDecodeObjectOwnershipBytesRejectsDuplicateKey(t *testing.T) {
	_, err := newTestCodec(t).DecodeObjectOwnership([]byte(`{"version":"v1","version":"v1"}`))

	requireErrorIs(t, err, ErrDuplicateKey)
}

// TestEncodeObjectOwnershipBytesUsesSharedWriter covers byte-slice ownership encoding.
func TestEncodeObjectOwnershipBytesUsesSharedWriter(t *testing.T) {
	data, err := newTestCodec(t).EncodeObjectOwnership(
		objectownership.Document{Version: objectownership.DocumentVersionV1},
	)
	requireNoError(t, err)

	if got, want := string(data), `{"version":"v1","desired":{"entries":[]}}`; got != want {
		t.Fatalf("encoded ownership = %s; want %s", got, want)
	}
}
