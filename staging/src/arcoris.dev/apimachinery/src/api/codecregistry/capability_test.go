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

var (
	_ codec.BaseCodec                  = fakeBaseCodec{}
	_ codec.ValueCodec                 = (*fakeValueByteCodec)(nil)
	_ codec.ObjectCodec                = (*fakeObjectByteCodec)(nil)
	_ codec.ObjectOwnershipCodec       = (*fakeOwnershipByteCodec)(nil)
	_ codec.Codec                      = (*fakeFullByteCodec)(nil)
	_ codec.ValueStreamCodec           = (*fakeValueStreamCodec)(nil)
	_ codec.ObjectStreamCodec          = (*fakeObjectStreamCodec)(nil)
	_ codec.ObjectOwnershipStreamCodec = (*fakeOwnershipStreamCodec)(nil)
	_ codec.StreamingCodec             = (*fakeFullStreamingCodec)(nil)
	_ codec.Codec                      = (*fakeByteAndStreamCodec)(nil)
	_ codec.StreamingCodec             = (*fakeByteAndStreamCodec)(nil)
)

func TestNewRejectsDeclaredTargetWithoutCapability(t *testing.T) {
	c := fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue),
	}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryError(t, err, "codecs[0].capabilities", ErrorReasonCapabilityMismatch)
}

func TestCapabilityMismatchDetailDeclaredTargetWithoutCapability(t *testing.T) {
	c := fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue),
	}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "declares target")
	requireRegistryDetailContains(t, err, "implements neither codec.ValueCodec nor codec.ValueStreamCodec")
}

func TestNewRejectsCapabilityWithoutDeclaredTarget(t *testing.T) {
	c := &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetObject),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryError(t, err, "codecs[0].capabilities", ErrorReasonCapabilityMismatch)
}

func TestCapabilityMismatchDetailCapabilityWithoutDeclaredTarget(t *testing.T) {
	c := &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetObject),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "does not declare target")
	requireRegistryDetailContains(t, err, "codec.ValueCodec or codec.ValueStreamCodec")
}

func TestCapabilityValidationTargetValueRequiresValueByteOrStream(t *testing.T) {
	tests := []struct {
		name  string
		codec codec.BaseCodec
	}{
		{name: "byte", codec: newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)},
		{name: "stream", codec: newValueStreamCodec(codec.FormatYAML, codec.MediaTypeYAML)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.codec)
			requireNoError(t, err)
		})
	}

	c := fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue),
	}
	_, err := New(c)
	requireErrorIs(t, err, ErrCapabilityMismatch)
}

func TestCapabilityValidationTargetObjectRequiresObjectByteOrStream(t *testing.T) {
	tests := []struct {
		name  string
		codec codec.BaseCodec
	}{
		{name: "byte", codec: newObjectByteCodec(codec.FormatJSON, codec.MediaTypeJSON)},
		{name: "stream", codec: newObjectStreamCodec(codec.FormatYAML, codec.MediaTypeYAML)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.codec)
			requireNoError(t, err)
		})
	}

	c := fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetObject),
	}
	_, err := New(c)
	requireErrorIs(t, err, ErrCapabilityMismatch)
}

func TestCapabilityValidationTargetObjectOwnershipRequiresOwnershipByteOrStream(t *testing.T) {
	tests := []struct {
		name  string
		codec codec.BaseCodec
	}{
		{name: "byte", codec: newOwnershipByteCodec(codec.FormatJSON, codec.MediaTypeJSON)},
		{name: "stream", codec: newOwnershipStreamCodec(codec.FormatYAML, codec.MediaTypeYAML)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.codec)
			requireNoError(t, err)
		})
	}

	c := fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetObjectOwnership),
	}
	_, err := New(c)
	requireErrorIs(t, err, ErrCapabilityMismatch)
}

func TestCapabilityValidationValueCapabilityRequiresDeclaredTarget(t *testing.T) {
	c := &fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetObject, codec.TargetObjectOwnership),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "codec.ValueCodec or codec.ValueStreamCodec")
}

func TestCapabilityValidationValueStreamCapabilityRequiresDeclaredTarget(t *testing.T) {
	c := &fakeFullStreamingCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetObject, codec.TargetObjectOwnership),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "codec.ValueCodec or codec.ValueStreamCodec")
}

func TestCapabilityValidationObjectCapabilityRequiresDeclaredTarget(t *testing.T) {
	c := &fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue, codec.TargetObjectOwnership),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "codec.ObjectCodec or codec.ObjectStreamCodec")
}

func TestCapabilityValidationObjectStreamCapabilityRequiresDeclaredTarget(t *testing.T) {
	c := &fakeFullStreamingCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue, codec.TargetObjectOwnership),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "codec.ObjectCodec or codec.ObjectStreamCodec")
}

func TestCapabilityValidationOwnershipCapabilityRequiresDeclaredTarget(t *testing.T) {
	c := &fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue, codec.TargetObject),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "codec.ObjectOwnershipCodec or codec.ObjectOwnershipStreamCodec")
}

func TestCapabilityValidationOwnershipStreamCapabilityRequiresDeclaredTarget(t *testing.T) {
	c := &fakeFullStreamingCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue, codec.TargetObject),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
	requireRegistryDetailContains(t, err, "codec.ObjectOwnershipCodec or codec.ObjectOwnershipStreamCodec")
}

func TestCapabilityValidationFullByteRequiresAllTargets(t *testing.T) {
	c := &fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue, codec.TargetObject),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
}

func TestCapabilityValidationFullStreamingRequiresAllTargets(t *testing.T) {
	c := &fakeFullStreamingCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(codec.FormatJSON, codec.MediaTypeJSON, codec.TargetValue, codec.TargetObject),
	}}

	_, err := New(c)

	requireErrorIs(t, err, ErrCapabilityMismatch)
}
