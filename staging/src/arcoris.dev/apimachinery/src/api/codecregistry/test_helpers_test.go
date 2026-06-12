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
	"strings"
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

func testRegistration(id string, c codec.BaseCodec) Registration {
	return Register(MustEntryID(id), c)
}

func testRegistry(t *testing.T, registrations ...Registration) Registry {
	t.Helper()

	registry, err := New(registrations...)
	requireNoError(t, err)

	return registry
}

func testValueByteRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newValueByteCodec(format, mediaType))
}

func testObjectByteRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newObjectByteCodec(format, mediaType))
}

func testOwnershipByteRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newOwnershipByteCodec(format, mediaType))
}

func testValueStreamRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newValueStreamCodec(format, mediaType))
}

func testObjectStreamRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newObjectStreamCodec(format, mediaType))
}

func testOwnershipStreamRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newOwnershipStreamCodec(format, mediaType))
}

func testFullByteRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newFullByteCodec(format, mediaType))
}

func testFullStreamRegistration(
	id string,
	format codec.Format,
	mediaType codec.MediaType,
) Registration {
	return testRegistration(id, newFullStreamingCodec(format, mediaType))
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

func (f fakeValueByteCodec) DecodeValue([]byte) (value.Value, error) {
	return value.NullValue(), nil
}

func (f fakeValueByteCodec) EncodeValue(value.Value) ([]byte, error) {
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

func (f fakeObjectByteCodec) DecodeObject([]byte) (codec.Object, error) {
	return codec.Object{}, nil
}

func (f fakeObjectByteCodec) EncodeObject(codec.Object) ([]byte, error) {
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
) (objectownership.State, error) {
	return objectownership.State{}, nil
}

func (f fakeOwnershipByteCodec) EncodeObjectOwnership(
	objectownership.State,
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

func (f fakeValueStreamCodec) DecodeValueFrom(io.Reader) (value.Value, error) {
	return value.NullValue(), nil
}

func (f fakeValueStreamCodec) EncodeValueTo(io.Writer, value.Value) error {
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

func (f fakeObjectStreamCodec) DecodeObjectFrom(io.Reader) (codec.Object, error) {
	return codec.Object{}, nil
}

func (f fakeObjectStreamCodec) EncodeObjectTo(io.Writer, codec.Object) error {
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
) (objectownership.State, error) {
	return objectownership.State{}, nil
}

func (f fakeOwnershipStreamCodec) EncodeObjectOwnershipTo(
	io.Writer,
	objectownership.State,
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

func (f fakeFullByteCodec) DecodeValue(data []byte) (value.Value, error) {
	return fakeValueByteCodec{}.DecodeValue(data)
}

func (f fakeFullByteCodec) EncodeValue(v value.Value) ([]byte, error) {
	return fakeValueByteCodec{}.EncodeValue(v)
}

func (f fakeFullByteCodec) DecodeObject(data []byte) (codec.Object, error) {
	return fakeObjectByteCodec{}.DecodeObject(data)
}

func (f fakeFullByteCodec) EncodeObject(obj codec.Object) ([]byte, error) {
	return fakeObjectByteCodec{}.EncodeObject(obj)
}

func (f fakeFullByteCodec) DecodeObjectOwnership(
	data []byte,
) (objectownership.State, error) {
	return fakeOwnershipByteCodec{}.DecodeObjectOwnership(data)
}

func (f fakeFullByteCodec) EncodeObjectOwnership(
	state objectownership.State,
) ([]byte, error) {
	return fakeOwnershipByteCodec{}.EncodeObjectOwnership(state)
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

func (f fakeFullStreamingCodec) DecodeValueFrom(r io.Reader) (value.Value, error) {
	return fakeValueStreamCodec{}.DecodeValueFrom(r)
}

func (f fakeFullStreamingCodec) EncodeValueTo(w io.Writer, v value.Value) error {
	return fakeValueStreamCodec{}.EncodeValueTo(w, v)
}

func (f fakeFullStreamingCodec) DecodeObjectFrom(r io.Reader) (codec.Object, error) {
	return fakeObjectStreamCodec{}.DecodeObjectFrom(r)
}

func (f fakeFullStreamingCodec) EncodeObjectTo(w io.Writer, obj codec.Object) error {
	return fakeObjectStreamCodec{}.EncodeObjectTo(w, obj)
}

func (f fakeFullStreamingCodec) DecodeObjectOwnershipFrom(
	r io.Reader,
) (objectownership.State, error) {
	return fakeOwnershipStreamCodec{}.DecodeObjectOwnershipFrom(r)
}

func (f fakeFullStreamingCodec) EncodeObjectOwnershipTo(
	w io.Writer,
	state objectownership.State,
) error {
	return fakeOwnershipStreamCodec{}.EncodeObjectOwnershipTo(w, state)
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

func (f fakeByteAndStreamCodec) DecodeValueFrom(r io.Reader) (value.Value, error) {
	return fakeFullStreamingCodec{}.DecodeValueFrom(r)
}

func (f fakeByteAndStreamCodec) EncodeValueTo(w io.Writer, v value.Value) error {
	return fakeFullStreamingCodec{}.EncodeValueTo(w, v)
}

func (f fakeByteAndStreamCodec) DecodeObjectFrom(r io.Reader) (codec.Object, error) {
	return fakeFullStreamingCodec{}.DecodeObjectFrom(r)
}

func (f fakeByteAndStreamCodec) EncodeObjectTo(w io.Writer, obj codec.Object) error {
	return fakeFullStreamingCodec{}.EncodeObjectTo(w, obj)
}

func (f fakeByteAndStreamCodec) DecodeObjectOwnershipFrom(
	r io.Reader,
) (objectownership.State, error) {
	return fakeFullStreamingCodec{}.DecodeObjectOwnershipFrom(r)
}

func (f fakeByteAndStreamCodec) EncodeObjectOwnershipTo(
	w io.Writer,
	state objectownership.State,
) error {
	return fakeFullStreamingCodec{}.EncodeObjectOwnershipTo(w, state)
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

func requireRegistryDetailContains(t *testing.T, err error, want string) {
	t.Helper()

	var registryError *Error
	if !errors.As(err, &registryError) {
		t.Fatalf("error = %T; want *Error", err)
	}
	if !strings.Contains(registryError.Detail, want) {
		t.Fatalf("detail = %q; want substring %q", registryError.Detail, want)
	}
}
