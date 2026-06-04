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

package codecregistry

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestIsNilCodecDetectsNilInterface(t *testing.T) {
	if !isNilCodec(nil) {
		t.Fatalf("isNilCodec(nil) = false")
	}
}

func TestIsNilCodecDetectsTypedNilPointer(t *testing.T) {
	var c *fakeValueByteCodec

	if !isNilCodec(c) {
		t.Fatalf("isNilCodec(typed nil pointer) = false")
	}
}

func TestIsNilCodecAcceptsConcreteValue(t *testing.T) {
	c := fakeBaseCodec{info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue)}

	if isNilCodec(c) {
		t.Fatalf("isNilCodec(concrete value) = true")
	}
}

func TestIsNilCodecAcceptsConcretePointer(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)

	if isNilCodec(c) {
		t.Fatalf("isNilCodec(pointer) = true")
	}
}
