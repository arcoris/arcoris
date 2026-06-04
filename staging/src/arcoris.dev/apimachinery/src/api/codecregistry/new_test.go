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

func TestNewEmpty(t *testing.T) {
	registry, err := New()
	requireNoError(t, err)

	if !registry.IsEmpty() || registry.Len() != 0 {
		t.Fatalf("registry empty = %v len = %d", registry.IsEmpty(), registry.Len())
	}
}

func TestZeroRegistryIsUsable(t *testing.T) {
	var registry Registry

	if !registry.IsEmpty() || registry.Len() != 0 {
		t.Fatalf("zero registry empty = %v len = %d", registry.IsEmpty(), registry.Len())
	}
	if _, ok := registry.LookupMediaType(codec.MediaTypeJSON); ok {
		t.Fatalf("zero registry lookup returned true")
	}
}

func TestNewRejectsNilCodec(t *testing.T) {
	_, err := New(nil)

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "codecs[0]", ErrorReasonInvalidCodec)
}

func TestNewRejectsTypedNilCodec(t *testing.T) {
	var c *fakeValueByteCodec

	_, err := New(c)

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "codecs[0]", ErrorReasonInvalidCodec)
}

func TestNewRejectsInvalidInfo(t *testing.T) {
	c := fakeBaseCodec{info: codec.Info{}}

	_, err := New(c)

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, codec.ErrInvalidInfo)
	requireRegistryError(t, err, "codecs[0].info", ErrorReasonInvalidInfo)
}

func TestNewReturnsZeroRegistryOnError(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireErrorIs(t, err, ErrDuplicateMediaType)
	if !registry.IsEmpty() {
		t.Fatalf("registry = %#v; want zero registry on error", registry)
	}
}

func TestNewCallsInfoOncePerCodec(t *testing.T) {
	c := &fakeCountingValueCodec{
		fakeValueByteCodec: newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	}

	_, err := New(c)
	requireNoError(t, err)

	if c.calls != 1 {
		t.Fatalf("Info calls = %d; want 1", c.calls)
	}
}

func TestNewAcceptsByteOnlyValueCodec(t *testing.T) {
	registry, err := New(newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValue(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupValue() = false")
	}
}

func TestNewAcceptsStreamOnlyValueCodec(t *testing.T) {
	registry, err := New(newValueStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupValueStream(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupValueStream() = false")
	}
}

func TestNewAcceptsFullByteCodec(t *testing.T) {
	registry, err := New(newFullByteCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupCodec(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupCodec() = false")
	}
}

func TestNewAcceptsFullStreamingCodec(t *testing.T) {
	registry, err := New(newFullStreamingCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupStreamingCodec(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupStreamingCodec() = false")
	}
}

func TestNewAcceptsCodecImplementingBothByteAndStream(t *testing.T) {
	registry, err := New(newByteAndStreamCodec(codec.FormatJSON, codec.MediaTypeJSON))
	requireNoError(t, err)

	if _, ok := registry.LookupCodec(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupCodec() = false")
	}
	if _, ok := registry.LookupStreamingCodec(codec.MediaTypeJSON); !ok {
		t.Fatalf("LookupStreamingCodec() = false")
	}
}
