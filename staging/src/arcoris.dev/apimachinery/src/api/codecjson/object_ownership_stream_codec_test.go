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
	"reflect"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/objectownership"
)

func TestDecodeObjectOwnershipFromMatchesDecodeObjectOwnership(t *testing.T) {
	c := newTestCodec(t)
	data := `{"desired":{"entries":[]}}`

	fromBytes, err := c.DecodeObjectOwnership([]byte(data))
	requireNoError(t, err)
	fromStream, err := c.DecodeObjectOwnershipFrom(strings.NewReader(data))
	requireNoError(t, err)

	if !reflect.DeepEqual(fromStream, fromBytes) {
		t.Fatalf("stream state = %#v; bytes state = %#v", fromStream, fromBytes)
	}
}

func TestEncodeObjectOwnershipToMatchesEncodeObjectOwnership(t *testing.T) {
	c := newTestCodec(t)
	state := objectownership.EmptyState()

	fromBytes, err := c.EncodeObjectOwnership(state)
	requireNoError(t, err)
	var buffer bytes.Buffer
	requireNoError(t, c.EncodeObjectOwnershipTo(&buffer, state))

	if buffer.String() != string(fromBytes) {
		t.Fatalf("stream = %s; bytes = %s", buffer.String(), fromBytes)
	}
}
