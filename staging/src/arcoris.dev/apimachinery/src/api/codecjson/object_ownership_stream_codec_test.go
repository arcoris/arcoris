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
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeObjectOwnershipFromMatchesDecodeObjectOwnership(t *testing.T) {
	c := newTestCodec(t)
	data := `{"version":"v1","desired":{"entries":[]}}`

	fromBytes, err := c.DecodeObjectOwnership([]byte(data))
	requireNoError(t, err)
	fromStream, err := c.DecodeObjectOwnershipFrom(strings.NewReader(data))
	requireNoError(t, err)

	if !reflect.DeepEqual(fromStream, fromBytes) {
		t.Fatalf("stream document = %#v; bytes document = %#v", fromStream, fromBytes)
	}
}

func TestEncodeObjectOwnershipToMatchesEncodeObjectOwnership(t *testing.T) {
	c := newTestCodec(t)
	doc := objectownership.Document{Version: objectownership.VersionV1}

	fromBytes, err := c.EncodeObjectOwnership(doc)
	requireNoError(t, err)
	var buffer bytes.Buffer
	requireNoError(t, c.EncodeObjectOwnershipTo(&buffer, doc))

	if buffer.String() != string(fromBytes) {
		t.Fatalf("stream = %s; bytes = %s", buffer.String(), fromBytes)
	}
}
