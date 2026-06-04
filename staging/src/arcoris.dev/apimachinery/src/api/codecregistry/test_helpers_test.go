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
	"errors"
	"io"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

func testInfo(format codec.Format, mediaType codec.MediaType, targets ...codec.Target) codec.Info {
	return codec.Info{
		Format:     format,
		MediaTypes: []codec.MediaType{mediaType},
		Targets:    targets,
	}
}

type fakeBaseCodec struct {
	info codec.Info
}

func (f fakeBaseCodec) Info() codec.Info {
	return f.info
}

type fakeValueByteCodec struct {
	fakeBaseCodec
}

type fakeCountingValueCodec struct {
	*fakeValueByteCodec
	calls int
}

func newValueByteCodec(format codec.Format, mediaType codec.MediaType) *fakeValueByteCodec {
	return &fakeValueByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaType, codec.TargetValue),
	}}
}

func (f *fakeCountingValueCodec) Info() codec.Info {
	f.calls++

	return f.fakeValueByteCodec.Info()
}

func (f fakeValueByteCodec) DecodeValue([]byte, codec.DecodeOptions) (value.Value, error) {
	return value.NullValue(), nil
}

func (f fakeValueByteCodec) EncodeValue(value.Value, codec.EncodeOptions) ([]byte, error) {
	return nil, nil
}

type fakeObjectByteCodec struct {
	fakeBaseCodec
}

func newObjectByteCodec(format codec.Format, mediaType codec.MediaType) *fakeObjectByteCodec {
	return &fakeObjectByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaType, codec.TargetObject),
	}}
}

func (f fakeObjectByteCodec) DecodeObject([]byte, codec.DecodeOptions) (codec.Object, error) {
	return codec.Object{}, nil
}

func (f fakeObjectByteCodec) EncodeObject(codec.Object, codec.EncodeOptions) ([]byte, error) {
	return nil, nil
}

type fakeOwnershipByteCodec struct {
	fakeBaseCodec
}

func newOwnershipByteCodec(format codec.Format, mediaType codec.MediaType) *fakeOwnershipByteCodec {
	return &fakeOwnershipByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaType, codec.TargetObjectOwnership),
	}}
}

func (f fakeOwnershipByteCodec) DecodeObjectOwnership(
	[]byte,
	codec.DecodeOptions,
) (objectownership.Document, error) {
	return objectownership.Document{}, nil
}

func (f fakeOwnershipByteCodec) EncodeObjectOwnership(
	objectownership.Document,
	codec.EncodeOptions,
) ([]byte, error) {
	return nil, nil
}

type fakeValueStreamCodec struct {
	fakeBaseCodec
}

func newValueStreamCodec(format codec.Format, mediaType codec.MediaType) *fakeValueStreamCodec {
	return &fakeValueStreamCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaType, codec.TargetValue),
	}}
}

func (f fakeValueStreamCodec) DecodeValueFrom(io.Reader, codec.DecodeOptions) (value.Value, error) {
	return value.NullValue(), nil
}

func (f fakeValueStreamCodec) EncodeValueTo(io.Writer, value.Value, codec.EncodeOptions) error {
	return nil
}

type fakeObjectStreamCodec struct {
	fakeBaseCodec
}

func newObjectStreamCodec(format codec.Format, mediaType codec.MediaType) *fakeObjectStreamCodec {
	return &fakeObjectStreamCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaType, codec.TargetObject),
	}}
}

func (f fakeObjectStreamCodec) DecodeObjectFrom(io.Reader, codec.DecodeOptions) (codec.Object, error) {
	return codec.Object{}, nil
}

func (f fakeObjectStreamCodec) EncodeObjectTo(io.Writer, codec.Object, codec.EncodeOptions) error {
	return nil
}

type fakeOwnershipStreamCodec struct {
	fakeBaseCodec
}

func newOwnershipStreamCodec(format codec.Format, mediaType codec.MediaType) *fakeOwnershipStreamCodec {
	return &fakeOwnershipStreamCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaType, codec.TargetObjectOwnership),
	}}
}

func (f fakeOwnershipStreamCodec) DecodeObjectOwnershipFrom(
	io.Reader,
	codec.DecodeOptions,
) (objectownership.Document, error) {
	return objectownership.Document{}, nil
}

func (f fakeOwnershipStreamCodec) EncodeObjectOwnershipTo(
	io.Writer,
	objectownership.Document,
	codec.EncodeOptions,
) error {
	return nil
}

type fakeFullByteCodec struct {
	fakeBaseCodec
}

func newFullByteCodec(format codec.Format, mediaType codec.MediaType) *fakeFullByteCodec {
	return &fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(
			format,
			mediaType,
			codec.TargetValue,
			codec.TargetObject,
			codec.TargetObjectOwnership,
		),
	}}
}

func (f fakeFullByteCodec) DecodeValue(data []byte, opts codec.DecodeOptions) (value.Value, error) {
	return fakeValueByteCodec{}.DecodeValue(data, opts)
}

func (f fakeFullByteCodec) EncodeValue(v value.Value, opts codec.EncodeOptions) ([]byte, error) {
	return fakeValueByteCodec{}.EncodeValue(v, opts)
}

func (f fakeFullByteCodec) DecodeObject(data []byte, opts codec.DecodeOptions) (codec.Object, error) {
	return fakeObjectByteCodec{}.DecodeObject(data, opts)
}

func (f fakeFullByteCodec) EncodeObject(obj codec.Object, opts codec.EncodeOptions) ([]byte, error) {
	return fakeObjectByteCodec{}.EncodeObject(obj, opts)
}

func (f fakeFullByteCodec) DecodeObjectOwnership(
	data []byte,
	opts codec.DecodeOptions,
) (objectownership.Document, error) {
	return fakeOwnershipByteCodec{}.DecodeObjectOwnership(data, opts)
}

func (f fakeFullByteCodec) EncodeObjectOwnership(
	doc objectownership.Document,
	opts codec.EncodeOptions,
) ([]byte, error) {
	return fakeOwnershipByteCodec{}.EncodeObjectOwnership(doc, opts)
}

type fakeFullStreamingCodec struct {
	fakeBaseCodec
}

func newFullStreamingCodec(format codec.Format, mediaType codec.MediaType) *fakeFullStreamingCodec {
	return &fakeFullStreamingCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(
			format,
			mediaType,
			codec.TargetValue,
			codec.TargetObject,
			codec.TargetObjectOwnership,
		),
	}}
}

func (f fakeFullStreamingCodec) DecodeValueFrom(r io.Reader, opts codec.DecodeOptions) (value.Value, error) {
	return fakeValueStreamCodec{}.DecodeValueFrom(r, opts)
}

func (f fakeFullStreamingCodec) EncodeValueTo(w io.Writer, v value.Value, opts codec.EncodeOptions) error {
	return fakeValueStreamCodec{}.EncodeValueTo(w, v, opts)
}

func (f fakeFullStreamingCodec) DecodeObjectFrom(r io.Reader, opts codec.DecodeOptions) (codec.Object, error) {
	return fakeObjectStreamCodec{}.DecodeObjectFrom(r, opts)
}

func (f fakeFullStreamingCodec) EncodeObjectTo(w io.Writer, obj codec.Object, opts codec.EncodeOptions) error {
	return fakeObjectStreamCodec{}.EncodeObjectTo(w, obj, opts)
}

func (f fakeFullStreamingCodec) DecodeObjectOwnershipFrom(
	r io.Reader,
	opts codec.DecodeOptions,
) (objectownership.Document, error) {
	return fakeOwnershipStreamCodec{}.DecodeObjectOwnershipFrom(r, opts)
}

func (f fakeFullStreamingCodec) EncodeObjectOwnershipTo(
	w io.Writer,
	doc objectownership.Document,
	opts codec.EncodeOptions,
) error {
	return fakeOwnershipStreamCodec{}.EncodeObjectOwnershipTo(w, doc, opts)
}

type fakeByteAndStreamCodec struct {
	fakeFullByteCodec
}

func newByteAndStreamCodec(format codec.Format, mediaType codec.MediaType) *fakeByteAndStreamCodec {
	return &fakeByteAndStreamCodec{fakeFullByteCodec: fakeFullByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(
			format,
			mediaType,
			codec.TargetValue,
			codec.TargetObject,
			codec.TargetObjectOwnership,
		),
	}}}
}

func (f fakeByteAndStreamCodec) DecodeValueFrom(r io.Reader, opts codec.DecodeOptions) (value.Value, error) {
	return fakeFullStreamingCodec{}.DecodeValueFrom(r, opts)
}

func (f fakeByteAndStreamCodec) EncodeValueTo(w io.Writer, v value.Value, opts codec.EncodeOptions) error {
	return fakeFullStreamingCodec{}.EncodeValueTo(w, v, opts)
}

func (f fakeByteAndStreamCodec) DecodeObjectFrom(r io.Reader, opts codec.DecodeOptions) (codec.Object, error) {
	return fakeFullStreamingCodec{}.DecodeObjectFrom(r, opts)
}

func (f fakeByteAndStreamCodec) EncodeObjectTo(w io.Writer, obj codec.Object, opts codec.EncodeOptions) error {
	return fakeFullStreamingCodec{}.EncodeObjectTo(w, obj, opts)
}

func (f fakeByteAndStreamCodec) DecodeObjectOwnershipFrom(
	r io.Reader,
	opts codec.DecodeOptions,
) (objectownership.Document, error) {
	return fakeFullStreamingCodec{}.DecodeObjectOwnershipFrom(r, opts)
}

func (f fakeByteAndStreamCodec) EncodeObjectOwnershipTo(
	w io.Writer,
	doc objectownership.Document,
	opts codec.EncodeOptions,
) error {
	return fakeFullStreamingCodec{}.EncodeObjectOwnershipTo(w, doc, opts)
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(..., %v)", err, target)
	}
}

func requireRegistryError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var registryError *Error
	if !errors.As(err, &registryError) {
		t.Fatalf("error = %T; want *Error", err)
	}
	if registryError.Path != path {
		t.Fatalf("path = %q; want %q", registryError.Path, path)
	}
	if registryError.Reason != reason {
		t.Fatalf("reason = %q; want %q", registryError.Reason, reason)
	}
}
